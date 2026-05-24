package indicators

import (
	"testing"
)

func TestBollingerBands_Update(t *testing.T) {
	bb := NewBollingerBands(5, 2)

	// All same prices — bands should narrow
	for i := 0; i < 10; i++ {
		band := bb.Update(100)
		if i >= 4 {
			if band.Upper < band.Middle || band.Lower > band.Middle {
				t.Errorf("invalid band at step %d: upper=%.2f middle=%.2f lower=%.2f",
					i, band.Upper, band.Middle, band.Lower)
			}
		}
	}
}

func TestBollingerBands_Ready(t *testing.T) {
	bb := NewBollingerBands(20, 2)
	for i := 0; i < 19; i++ {
		bb.Update(100)
		if bb.Ready() {
			t.Errorf("should not be ready at step %d", i+1)
		}
	}
	bb.Update(100)
	if !bb.Ready() {
		t.Errorf("should be ready after 20 updates")
	}
}

func TestPercentB(t *testing.T) {
	band := BollingerBand{Upper: 110, Middle: 100, Lower: 90}

	if band.PercentB(100) != 0.5 {
		t.Errorf("expected 0.5 at middle, got %.2f", band.PercentB(100))
	}

	if band.PercentB(90) != 0.0 {
		t.Errorf("expected 0.0 at lower, got %.2f", band.PercentB(90))
	}

	if band.PercentB(110) != 1.0 {
		t.Errorf("expected 1.0 at upper, got %.2f", band.PercentB(110))
	}
}

func TestBandwidth(t *testing.T) {
	band := BollingerBand{Upper: 110, Middle: 100, Lower: 90}

	bw := band.Bandwidth()
	if bw != 0.2 {
		t.Errorf("expected 0.2, got %.2f", bw)
	}
}

func TestBollingerValues(t *testing.T) {
	prices := make([]float64, 50)
	for i := range prices {
		prices[i] = 100 + float64(i)*0.2
	}

	result := BollingerValues(prices, 20, 2)
	if len(result) != len(prices) {
		t.Errorf("expected %d results, got %d", len(prices), len(result))
	}
}
