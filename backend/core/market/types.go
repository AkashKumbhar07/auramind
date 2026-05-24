package market

type Symbol struct {
	Base  string
	Quote string
}

func (s Symbol) Pair() string {
	return s.Base + s.Quote
}

type Side string

const (
	SideBuy  Side = "buy"
	SideSell Side = "sell"
)

type Tick struct {
	Symbol    string  `json:"symbol"`
	Exchange  string  `json:"exchange"`
	Price     float64 `json:"price"`
	Quantity  float64 `json:"quantity"`
	Side      Side    `json:"side"`
	Timestamp int64   `json:"timestamp"`
}

type Trade struct {
	Symbol    string  `json:"symbol"`
	Exchange  string  `json:"exchange"`
	Price     float64 `json:"price"`
	Quantity  float64 `json:"quantity"`
	Side      Side    `json:"side"`
	TradeID   int64   `json:"trade_id"`
	Timestamp int64   `json:"timestamp"`
}

type OHLCV struct {
	Symbol    string  `json:"symbol"`
	Exchange  string  `json:"exchange"`
	Interval  string  `json:"interval"`
	Open      float64 `json:"open"`
	High      float64 `json:"high"`
	Low       float64 `json:"low"`
	Close     float64 `json:"close"`
	Volume    float64 `json:"volume"`
	Timestamp int64   `json:"timestamp"`
	Closed    bool    `json:"closed"`
}

func (c *OHLCV) Update(price float64) {
	if price > c.High {
		c.High = price
	}
	if price < c.Low || c.Low == 0 {
		c.Low = price
	}
	c.Close = price
}

type OrderBookLevel struct {
	Price    float64 `json:"price"`
	Quantity float64 `json:"quantity"`
}

type OrderBook struct {
	Symbol    string            `json:"symbol"`
	Exchange  string            `json:"exchange"`
	Bids      []OrderBookLevel  `json:"bids"`
	Asks      []OrderBookLevel  `json:"asks"`
	Timestamp int64             `json:"timestamp"`
}

type Ticker struct {
	Symbol    string  `json:"symbol"`
	Exchange  string  `json:"exchange"`
	LastPrice float64 `json:"last_price"`
	Volume24h float64 `json:"volume_24h"`
	High24h   float64 `json:"high_24h"`
	Low24h    float64 `json:"low_24h"`
	Change24h float64 `json:"change_24h"`
	Timestamp int64   `json:"timestamp"`
}

type EventType string

const (
	EventTrade     EventType = "trade"
	EventTicker    EventType = "ticker"
	EventOHLCV     EventType = "ohlcv"
	EventOrderBook EventType = "orderbook"
)

type MarketEvent struct {
	Type      EventType `json:"type"`
	Symbol    string    `json:"symbol"`
	Exchange  string    `json:"exchange"`
	Data      any       `json:"data"`
	Timestamp int64     `json:"timestamp"`
}

type CandleInterval string

const (
	Interval1m  CandleInterval = "1m"
	Interval5m  CandleInterval = "5m"
	Interval15m CandleInterval = "15m"
	Interval1h  CandleInterval = "1h"
	Interval4h  CandleInterval = "4h"
	Interval1d  CandleInterval = "1d"
)
