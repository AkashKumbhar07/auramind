package indicators

import (
	"testing"
)

func TestMACD_Update(t *testing.T) {
	macd := NewMACD(12, 26, 9)

	for i := 0; i < 30; i++ {
		line := macd.Update(100 + float64(i))
		if i > 26 && line.Histogram == 0 {
			t.Errorf("expected non-zero histogram at step %d", i)
		}
	}
}

func TestMACD_Ready(t *testing.T) {
	macd := NewMACD(12, 26, 9)
	for i := 0; i < 25; i++ {
		macd.Update(100)
		if macd.Ready() {
			t.Errorf("should not be ready at step %d", i+1)
		}
	}
	macd.Update(100)
	if !macd.Ready() {
		t.Errorf("should be ready after 27 updates")
	}
}

func TestMACDValues(t *testing.T) {
	prices := make([]float64, 50)
	for i := range prices {
		prices[i] = 100 + float64(i)*0.5
	}

	result := MACDValues(prices, 12, 26, 9)
	if len(result) != len(prices) {
		t.Errorf("expected %d results, got %d", len(prices), len(result))
	}
}
