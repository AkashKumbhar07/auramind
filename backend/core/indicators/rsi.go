package indicators

type RSI struct {
	period  int
	gain    float64
	loss    float64
	prev    float64
	ready   bool
	count   int
	average bool
	value   float64
}

func NewRSI(period int) *RSI {
	return &RSI{
		period: period,
	}
}

func (r *RSI) Update(price float64) float64 {
	r.count++

	if !r.ready {
		r.prev = price
		r.ready = true
		return 50
	}

	diff := price - r.prev
	r.prev = price

	var gain, loss float64
	if diff > 0 {
		gain = diff
	} else {
		loss = -diff
	}

	if !r.average {
		r.gain += gain
		r.loss += loss

		if r.count == r.period+1 {
			r.gain /= float64(r.period)
			r.loss /= float64(r.period)
			r.average = true
			r.value = r.calcRSI()
		}
		return 50
	}

	r.gain = (r.gain*float64(r.period-1) + gain) / float64(r.period)
	r.loss = (r.loss*float64(r.period-1) + loss) / float64(r.period)
	r.value = r.calcRSI()

	return r.value
}

func (r *RSI) calcRSI() float64 {
	if r.loss == 0 {
		return 100
	}
	rs := r.gain / r.loss
	return 100 - (100 / (1 + rs))
}

func (r *RSI) Value() float64 {
	return r.value
}

func (r *RSI) Ready() bool {
	return r.count > r.period
}

func (r *RSI) Reset() {
	r.gain = 0
	r.loss = 0
	r.prev = 0
	r.ready = false
	r.count = 0
	r.average = false
	r.value = 50
}

func RSIValues(prices []float64, period int) []float64 {
	if len(prices) == 0 {
		return nil
	}

	result := make([]float64, len(prices))
	rsi := NewRSI(period)

	for i, p := range prices {
		result[i] = rsi.Update(p)
	}

	return result
}
