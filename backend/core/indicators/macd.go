package indicators

type MACDLine struct {
	Value     float64 `json:"value"`
	Signal    float64 `json:"signal"`
	Histogram float64 `json:"histogram"`
}

type MACD struct {
	fast   *EMA
	slow   *EMA
	signal *EMA
}

func NewMACD(fastPeriod, slowPeriod, signalPeriod int) *MACD {
	return &MACD{
		fast:   NewEMA(fastPeriod),
		slow:   NewEMA(slowPeriod),
		signal: NewEMA(signalPeriod),
	}
}

func (m *MACD) Update(price float64) MACDLine {
	fastVal := m.fast.Update(price)
	slowVal := m.slow.Update(price)
	macdLine := fastVal - slowVal
	signalVal := m.signal.Update(macdLine)

	return MACDLine{
		Value:     macdLine,
		Signal:    signalVal,
		Histogram: macdLine - signalVal,
	}
}

func (m *MACD) Ready() bool {
	return m.slow.Ready() && m.signal.Ready()
}

func (m *MACD) Reset() {
	m.fast.Reset()
	m.slow.Reset()
	m.signal.Reset()
}

func MACDValues(prices []float64, fast, slow, signal int) []MACDLine {
	if len(prices) == 0 {
		return nil
	}

	result := make([]MACDLine, len(prices))
	macd := NewMACD(fast, slow, signal)

	for i, p := range prices {
		result[i] = macd.Update(p)
	}

	return result
}
