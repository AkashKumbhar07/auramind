package indicators

import (
	"math"
	"testing"
)

func TestEMA_Update(t *testing.T) {
	ema := NewEMA(3)

	values := []float64{10, 12, 11, 13, 14, 15}
	expected := []float64{10, 11, 11, 12, 13, 14}

	for i, v := range values {
		result := ema.Update(v)
		rounded := math.Round(result)
		if rounded != expected[i] {
			t.Errorf("step %d: expected %.0f, got %.0f", i, expected[i], rounded)
		}
	}
}

func TestEMA_Ready(t *testing.T) {
	ema := NewEMA(5)
	for i := 0; i < 4; i++ {
		ema.Update(100)
		if ema.Ready() {
			t.Errorf("should not be ready at count %d", i+1)
		}
	}
	ema.Update(100)
	if !ema.Ready() {
		t.Errorf("should be ready after 5 updates")
	}
}

func TestEMAValues(t *testing.T) {
	prices := []float64{10, 12, 11, 13, 14, 15, 14, 16, 17, 18}
	result := EMAValues(prices, 5)

	if len(result) != len(prices) {
		t.Errorf("expected %d results, got %d", len(prices), len(result))
	}

	if result[0] != prices[0] {
		t.Errorf("first value should equal first price")
	}
}
