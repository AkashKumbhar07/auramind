package indicators

import (
	"math"
	"testing"
)

func TestRSI_Update(t *testing.T) {
	rsi := NewRSI(14)

	// Up trend — RSI should go above 70
	for i := 0; i < 20; i++ {
		rsi.Update(100 + float64(i))
	}

	if rsi.Value() < 70 {
		t.Errorf("expected RSI > 70 in uptrend, got %.2f", rsi.Value())
	}
}

func TestRSI_Oversold(t *testing.T) {
	rsi := NewRSI(14)

	// Down trend — RSI should go below 30
	for i := 0; i < 20; i++ {
		rsi.Update(100 - float64(i)*2)
	}

	if rsi.Value() > 30 {
		t.Errorf("expected RSI < 30 in downtrend, got %.2f", rsi.Value())
	}
}

func TestRSI_Ready(t *testing.T) {
	rsi := NewRSI(5)
	for i := 0; i < 5; i++ {
		rsi.Update(100)
		if rsi.Ready() {
			t.Errorf("should not be ready at count %d", i+1)
		}
	}
	rsi.Update(100)
	if !rsi.Ready() {
		t.Errorf("should be ready after 6 updates")
	}
}

func TestRSIValues(t *testing.T) {
	prices := make([]float64, 30)
	for i := range prices {
		prices[i] = 100 + float64(i)
	}

	result := RSIValues(prices, 14)
	if len(result) != len(prices) {
		t.Errorf("expected %d results, got %d", len(prices), len(result))
	}

	last := result[len(result)-1]
	if math.IsNaN(last) {
		t.Errorf("last RSI value is NaN")
	}
}
