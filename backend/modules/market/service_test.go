package market

import (
	"context"
	"testing"

	"go.uber.org/zap"
)

type mockRepo struct{}

func (m *mockRepo) GetPairs(ctx context.Context, exchange string) ([]MarketPair, error) {
	return []MarketPair{
		{Symbol: "BTCUSDT", Exchange: "binance", Active: true},
	}, nil
}

func (m *mockRepo) GetCandles(ctx context.Context, symbol, exchange, interval string, limit int) ([]Candle, error) {
	candles := make([]Candle, 100)
	for i := range candles {
		candles[i] = Candle{
			Symbol: symbol, Exchange: exchange, Interval: interval,
			Close: 50000 + float64(i),
		}
	}
	return candles, nil
}

func (m *mockRepo) SaveCandle(ctx context.Context, candle *Candle) error { return nil }
func (m *mockRepo) BulkSaveCandles(ctx context.Context, candles []*Candle) error { return nil }

func TestService_Analyze(t *testing.T) {
	svc := NewService(&mockRepo{}, zap.NewNop())

	candles := make([]Candle, 100)
	for i := range candles {
		candles[i] = Candle{Symbol: "BTCUSDT", Close: 50000 + float64(i)}
	}

	resp, err := svc.Analyze(context.Background(), &AnalyzeRequest{
		Symbol:  "BTCUSDT",
		Candles: candles,
	})
	if err != nil {
		t.Fatalf("analyze failed: %v", err)
	}

	if resp.Symbol != "BTCUSDT" {
		t.Errorf("expected BTCUSDT, got %s", resp.Symbol)
	}

	if resp.BullishProbability < 0 || resp.BullishProbability > 100 {
		t.Errorf("probability out of range: %.2f", resp.BullishProbability)
	}

	if resp.Indicators.RSI == nil {
		t.Error("expected RSI indicator")
	}

	if len(resp.Reasons) == 0 {
		t.Error("expected at least one reason")
	}
}

func TestService_Analyze_InsufficientData(t *testing.T) {
	svc := NewService(&mockRepo{}, zap.NewNop())

	_, err := svc.Analyze(context.Background(), &AnalyzeRequest{
		Symbol:  "BTCUSDT",
		Candles: []Candle{{Close: 100}},
	})
	if err == nil {
		t.Error("expected error for insufficient data")
	}
}

func TestCalcStopLoss(t *testing.T) {
	svc := NewService(&mockRepo{}, zap.NewNop())

	candles := make([]Candle, 20)
	for i := range candles {
		candles[i] = Candle{Close: 100}
	}

	sl := svc.calcStopLoss(candles)
	if sl <= 0 {
		t.Errorf("expected positive stop loss, got %.2f", sl)
	}
}

func TestCalcRiskLevel(t *testing.T) {
	svc := NewService(&mockRepo{}, zap.NewNop())

	tests := []struct {
		score    float64
		expected string
	}{
		{50, "low"},
		{20, "medium"},  // extreme score but no volatility
		{85, "medium"},  // extreme score but no volatility
	}

	for _, tt := range tests {
		risk := svc.calcRiskLevel(tt.score, IndicatorResult{})
		if risk != tt.expected {
			t.Errorf("score %.0f: expected %s, got %s", tt.score, tt.expected, risk)
		}
	}
}
