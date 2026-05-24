package market

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"go.uber.org/zap"
)

type ExchangeClient interface {
	Connect(ctx context.Context) error
	Subscribe(symbols ...string) error
	Trades() <-chan Trade
	Tickers() <-chan Ticker
	OHLCV() <-chan OHLCV
	Errors() <-chan error
	Close() error
}

type Exchange string

const (
	ExchangeBinance Exchange = "binance"
	ExchangeBybit   Exchange = "bybit"
)

type IngestionConfig struct {
	Exchange    Exchange
	Symbols     []string
	WSURL       string
	RestURL     string
	ReconnectMs int
	Logger      *zap.Logger
}

type IngestionService struct {
	cfg     IngestionConfig
	client  ExchangeClient
	trades  chan Trade
	tickers chan Ticker
	ohlcvs  chan OHLCV
	events  chan MarketEvent
	errs    chan error
	mu      sync.Mutex
	running bool
	done    chan struct{}
}

func NewIngestion(cfg IngestionConfig) *IngestionService {
	return &IngestionService{
		cfg:     cfg,
		trades:  make(chan Trade, 1024),
		tickers: make(chan Ticker, 256),
		ohlcvs:  make(chan OHLCV, 512),
		events:  make(chan MarketEvent, 1024),
		errs:    make(chan error, 16),
		done:    make(chan struct{}),
	}
}

func (s *IngestionService) Start(ctx context.Context) error {
	s.mu.Lock()
	if s.running {
		s.mu.Unlock()
		return nil
	}
	s.running = true
	s.mu.Unlock()

	go s.run(ctx)
	return nil
}

func (s *IngestionService) run(ctx context.Context) {
	defer close(s.done)

	reconnectDelay := time.Duration(s.cfg.ReconnectMs) * time.Millisecond
	if reconnectDelay == 0 {
		reconnectDelay = 5 * time.Second
	}

	for {
		select {
		case <-ctx.Done():
			s.cfg.Logger.Info("ingestion stopped", zap.String("exchange", string(s.cfg.Exchange)))
			return
		default:
		}

		s.cfg.Logger.Info("connecting to exchange",
			zap.String("exchange", string(s.cfg.Exchange)),
			zap.String("ws_url", s.cfg.WSURL),
		)

		err := s.connectAndStream(ctx)
		if err != nil {
			s.cfg.Logger.Error("stream error, reconnecting",
				zap.String("exchange", string(s.cfg.Exchange)),
				zap.Error(err),
				zap.Duration("delay", reconnectDelay),
			)

			select {
			case s.errs <- err:
			default:
			}

			select {
			case <-ctx.Done():
				return
			case <-time.After(reconnectDelay):
			}
		}
	}
}

func (s *IngestionService) connectAndStream(ctx context.Context) error {
	client, err := NewBinanceClient(s.cfg.WSURL, s.cfg.Logger)
	if err != nil {
		return fmt.Errorf("create client: %w", err)
	}
	defer client.Close()

	if err := client.Connect(ctx); err != nil {
		return fmt.Errorf("connect: %w", err)
	}

	if err := client.Subscribe(s.cfg.Symbols...); err != nil {
		return fmt.Errorf("subscribe: %w", err)
	}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()

		case trade, ok := <-client.Trades():
			if !ok {
				return fmt.Errorf("trade channel closed")
			}
			s.trades <- trade
			s.events <- MarketEvent{
				Type: EventTrade, Symbol: trade.Symbol,
				Exchange: string(s.cfg.Exchange), Data: trade, Timestamp: trade.Timestamp,
			}

		case ticker, ok := <-client.Tickers():
			if !ok {
				return fmt.Errorf("ticker channel closed")
			}
			s.tickers <- ticker
			s.events <- MarketEvent{
				Type: EventTicker, Symbol: ticker.Symbol,
				Exchange: string(s.cfg.Exchange), Data: ticker, Timestamp: ticker.Timestamp,
			}

		case ohlcv, ok := <-client.OHLCV():
			if !ok {
				return fmt.Errorf("ohlcv channel closed")
			}
			s.ohlcvs <- ohlcv
			s.events <- MarketEvent{
				Type: EventOHLCV, Symbol: ohlcv.Symbol,
				Exchange: string(s.cfg.Exchange), Data: ohlcv, Timestamp: ohlcv.Timestamp,
			}

		case err, ok := <-client.Errors():
			if !ok {
				return fmt.Errorf("error channel closed")
			}
			return err
		}
	}
}

func (s *IngestionService) Stop() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.running = false
}

func (s *IngestionService) Trades() <-chan Trade {
	return s.trades
}

func (s *IngestionService) Tickers() <-chan Ticker {
	return s.tickers
}

func (s *IngestionService) OHLCV() <-chan OHLCV {
	return s.ohlcvs
}

func (s *IngestionService) Events() <-chan MarketEvent {
	return s.events
}

func (s *IngestionService) Errors() <-chan error {
	return s.errs
}

// BinanceWebSocket client implementation

type binanceTrade struct {
	EventType string `json:"e"`
	TradeID   int64  `json:"t"`
	Symbol    string `json:"s"`
	Price     string `json:"p"`
	Quantity  string `json:"q"`
	Time      int64  `json:"T"`
	IsBuyer   bool   `json:"m"`
}

type binanceTicker struct {
	EventType  string `json:"e"`
	Symbol     string `json:"s"`
	Price      string `json:"c"`
	Volume     string `json:"v"`
	High       string `json:"h"`
	Low        string `json:"l"`
	Change     string `json:"p"`
	Time       int64  `json:"E"`
}

type binanceKline struct {
	EventType string        `json:"e"`
	Symbol    string        `json:"s"`
	Kline     binanceKlineData `json:"k"`
}

type binanceKlineData struct {
	Start  int64  `json:"t"`
	End    int64  `json:"T"`
	Symbol string `json:"s"`
	Open   string `json:"o"`
	High   string `json:"h"`
	Low    string `json:"l"`
	Close  string `json:"c"`
	Volume string `json:"v"`
	Closed bool   `json:"x"`
}

type binanceClient struct {
	wsURL     string
	conn      *WebSocketConn
	trades    chan Trade
	tickers   chan Ticker
	ohlcvs    chan OHLCV
	errs      chan error
	logger    *zap.Logger
	subMu     sync.Mutex
	subscribed map[string]bool
}

func NewBinanceClient(wsURL string, logger *zap.Logger) (*binanceClient, error) {
	return &binanceClient{
		wsURL:      wsURL,
		trades:     make(chan Trade, 1024),
		tickers:    make(chan Ticker, 256),
		ohlcvs:     make(chan OHLCV, 512),
		errs:       make(chan error, 16),
		logger:     logger,
		subscribed: make(map[string]bool),
	}, nil
}

func (c *binanceClient) Connect(ctx context.Context) error {
	conn, err := Dial(ctx, c.wsURL)
	if err != nil {
		return err
	}
	c.conn = conn
	go c.readLoop()
	return nil
}

func (c *binanceClient) readLoop() {
	defer func() {
		close(c.trades)
		close(c.tickers)
		close(c.ohlcvs)
		close(c.errs)
	}()

	for {
		msg, err := c.conn.Read()
		if err != nil {
			c.errs <- err
			return
		}

		c.dispatch(msg)
	}
}

func (c *binanceClient) dispatch(msg []byte) {
	raw := string(msg)

	switch {
	case strings.Contains(raw, `"e":"trade"`):
		c.handleTrade(raw)
	case strings.Contains(raw, `"e":"24hrTicker"`):
		c.handleTicker(raw)
	case strings.Contains(raw, `"e":"kline"`):
		c.handleKline(raw)
	}
}

func parseFloat(s string) float64 {
	v, _ := strconv.ParseFloat(s, 64)
	return v
}

func (c *binanceClient) handleTrade(raw string) {
	var t binanceTrade
	if err := jsonUnmarshal([]byte(raw), &t); err != nil {
		c.errs <- fmt.Errorf("parse trade: %w", err)
		return
	}

	side := SideBuy
	if !t.IsBuyer {
		side = SideSell
	}

	c.trades <- Trade{
		Symbol:    strings.ToUpper(t.Symbol),
		Exchange:  "binance",
		Price:     parseFloat(t.Price),
		Quantity:  parseFloat(t.Quantity),
		Side:      side,
		TradeID:   t.TradeID,
		Timestamp: t.Time,
	}
}

func (c *binanceClient) handleTicker(raw string) {
	var t binanceTicker
	if err := jsonUnmarshal([]byte(raw), &t); err != nil {
		c.errs <- fmt.Errorf("parse ticker: %w", err)
		return
	}

	c.tickers <- Ticker{
		Symbol:    strings.ToUpper(t.Symbol),
		Exchange:  "binance",
		LastPrice: parseFloat(t.Price),
		Volume24h: parseFloat(t.Volume),
		High24h:   parseFloat(t.High),
		Low24h:    parseFloat(t.Low),
		Change24h: parseFloat(t.Change),
		Timestamp: t.Time,
	}
}

func (c *binanceClient) handleKline(raw string) {
	var k binanceKline
	if err := jsonUnmarshal([]byte(raw), &k); err != nil {
		c.errs <- fmt.Errorf("parse kline: %w", err)
		return
	}

	c.ohlcvs <- OHLCV{
		Symbol:    strings.ToUpper(k.Kline.Symbol),
		Exchange:  "binance",
		Interval:  intervalMap(k.Kline.End - k.Kline.Start),
		Open:      parseFloat(k.Kline.Open),
		High:      parseFloat(k.Kline.High),
		Low:       parseFloat(k.Kline.Low),
		Close:     parseFloat(k.Kline.Close),
		Volume:    parseFloat(k.Kline.Volume),
		Timestamp: k.Kline.End,
		Closed:    k.Kline.Closed,
	}
}

func intervalMap(diffMs int64) string {
	switch {
	case diffMs <= 60000:
		return "1m"
	case diffMs <= 300000:
		return "5m"
	case diffMs <= 900000:
		return "15m"
	case diffMs <= 3600000:
		return "1h"
	case diffMs <= 14400000:
		return "4h"
	default:
		return "1d"
	}
}

func (c *binanceClient) Subscribe(symbols ...string) error {
	c.subMu.Lock()
	defer c.subMu.Unlock()

	var params []string
	for _, s := range symbols {
		lower := strings.ToLower(s)
		if c.subscribed[lower] {
			continue
		}
		c.subscribed[lower] = true
		params = append(params,
			fmt.Sprintf("%s@trade", lower),
			fmt.Sprintf("%s@ticker", lower),
			fmt.Sprintf("%s@kline_1m", lower),
		)
	}

	if len(params) == 0 {
		return nil
	}

	subMsg := map[string]any{
		"method": "SUBSCRIBE",
		"params": params,
		"id":     1,
	}

	return c.conn.WriteJSON(subMsg)
}

func (c *binanceClient) Trades() <-chan Trade {
	return c.trades
}

func (c *binanceClient) Tickers() <-chan Ticker {
	return c.tickers
}

func (c *binanceClient) OHLCV() <-chan OHLCV {
	return c.ohlcvs
}

func (c *binanceClient) Errors() <-chan error {
	return c.errs
}

func (c *binanceClient) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}
