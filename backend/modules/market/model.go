package market

type MarketPair struct {
	ID        string `json:"id"`
	Symbol    string `json:"symbol"`
	BaseAsset string `json:"base_asset"`
	QuoteAsset string `json:"quote_asset"`
	Exchange  string `json:"exchange"`
	Active    bool   `json:"active"`
}

type Candle struct {
	Symbol    string  `json:"symbol"`
	Exchange  string  `json:"exchange"`
	Interval  string  `json:"interval"`
	Open      float64 `json:"open"`
	High      float64 `json:"high"`
	Low       float64 `json:"low"`
	Close     float64 `json:"close"`
	Volume    float64 `json:"volume"`
	Timestamp int64   `json:"timestamp"`
}

type IndicatorResult struct {
	Symbol string `json:"symbol"`
	RSI    *float64 `json:"rsi,omitempty"`
	MACD   *MACDValue `json:"macd,omitempty"`
	BB     *BBValue   `json:"bollinger,omitempty"`
	EMA9   *float64   `json:"ema_9,omitempty"`
	EMA21  *float64   `json:"ema_21,omitempty"`
}

type MACDValue struct {
	Value     float64 `json:"value"`
	Signal    float64 `json:"signal"`
	Histogram float64 `json:"histogram"`
}

type BBValue struct {
	Upper  float64 `json:"upper"`
	Middle float64 `json:"middle"`
	Lower  float64 `json:"lower"`
}
