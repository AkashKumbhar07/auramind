package indicators

import "math"

type BollingerBand struct {
	Upper  float64 `json:"upper"`
	Middle float64 `json:"middle"`
	Lower  float64 `json:"lower"`
	Width  float64 `json:"width"`
}

type BollingerBands struct {
	period   int
	stdDev   float64
	prices   []float64
	sum      float64
	sumSq    float64
	value    BollingerBand
	ready    bool
}

func NewBollingerBands(period int, stdDev float64) *BollingerBands {
	return &BollingerBands{
		period: period,
		stdDev: stdDev,
		prices: make([]float64, 0, period),
	}
}

func (b *BollingerBands) Update(price float64) BollingerBand {
	if len(b.prices) >= b.period {
		oldest := b.prices[0]
		b.sum -= oldest
		b.sumSq -= oldest * oldest
		b.prices = b.prices[1:]
	}

	b.prices = append(b.prices, price)
	b.sum += price
	b.sumSq += price * price

	if len(b.prices) < b.period {
		return BollingerBand{}
	}

	b.ready = true
	mean := b.sum / float64(len(b.prices))
	variance := (b.sumSq / float64(len(b.prices))) - (mean * mean)
	if variance < 0 {
		variance = 0
	}
	std := math.Sqrt(variance)

	band := std * b.stdDev

	b.value = BollingerBand{
		Middle: mean,
		Upper:  mean + band,
		Lower:  mean - band,
		Width:  band * 2,
	}

	return b.value
}

func (b *BollingerBands) Value() BollingerBand {
	return b.value
}

func (b *BollingerBands) Ready() bool {
	return b.ready
}

func (b *BollingerBands) Reset() {
	b.prices = make([]float64, 0, b.period)
	b.sum = 0
	b.sumSq = 0
	b.ready = false
}

func BollingerValues(prices []float64, period int, stdDev float64) []BollingerBand {
	if len(prices) == 0 {
		return nil
	}

	result := make([]BollingerBand, len(prices))
	bb := NewBollingerBands(period, stdDev)

	for i, p := range prices {
		result[i] = bb.Update(p)
	}

	return result
}

func (b *BollingerBand) PercentB(price float64) float64 {
	if b.Upper == b.Lower {
		return 0.5
	}
	return (price - b.Lower) / (b.Upper - b.Lower)
}

func (b *BollingerBand) Bandwidth() float64 {
	if b.Middle == 0 {
		return 0
	}
	return (b.Upper - b.Lower) / b.Middle
}
