package indicators

type EMA struct {
	period   int
	multiplier float64
	value    float64
	ready    bool
	count    int
}

func NewEMA(period int) *EMA {
	return &EMA{
		period:     period,
		multiplier: 2.0 / float64(period+1),
	}
}

func (e *EMA) Update(price float64) float64 {
	e.count++

	if !e.ready {
		e.value = price
		e.ready = true
		return e.value
	}

	e.value = (price-e.value)*e.multiplier + e.value
	return e.value
}

func (e *EMA) Value() float64 {
	return e.value
}

func (e *EMA) Ready() bool {
	return e.ready && e.count >= e.period
}

func (e *EMA) Reset() {
	e.value = 0
	e.ready = false
	e.count = 0
}

func EMAValues(prices []float64, period int) []float64 {
	if len(prices) == 0 {
		return nil
	}

	result := make([]float64, len(prices))
	ema := NewEMA(period)

	for i, p := range prices {
		result[i] = ema.Update(p)
	}

	return result
}
