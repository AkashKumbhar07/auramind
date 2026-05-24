package market

type GetPairsRequest struct {
	Exchange string `json:"exchange"`
}

type GetPairsResponse struct {
	Pairs []MarketPair `json:"pairs"`
}

type GetCandlesRequest struct {
	Symbol   string `json:"symbol" validate:"required"`
	Exchange string `json:"exchange"`
	Interval string `json:"interval" validate:"required,oneof=1m 5m 15m 1h 4h 1d"`
	Limit    int    `json:"limit"`
}

type GetCandlesResponse struct {
	Candles []Candle `json:"candles"`
}

type GetTickerRequest struct {
	Symbol   string `json:"symbol" validate:"required"`
	Exchange string `json:"exchange"`
}

type GetTickerResponse struct {
	Symbol    string  `json:"symbol"`
	LastPrice float64 `json:"last_price"`
	Volume24h float64 `json:"volume_24h"`
	High24h   float64 `json:"high_24h"`
	Low24h    float64 `json:"low_24h"`
	Change24h float64 `json:"change_24h"`
}

type AnalyzeRequest struct {
	Symbol string  `json:"symbol" validate:"required"`
	Candles []Candle `json:"candles" validate:"required,min=30"`
}

type AnalyzeResponse struct {
	Symbol             string          `json:"symbol"`
	BullishProbability float64         `json:"bullish_probability"`
	ConfidenceScore    float64         `json:"confidence_score"`
	RiskLevel          string          `json:"risk_level"`
	SuggestedStopLoss  float64         `json:"suggested_stop_loss"`
	Indicators         IndicatorResult `json:"indicators"`
	Reasons            []string        `json:"reasons"`
	Invalidation       string          `json:"invalidation_condition"`
}

type SubscribeRequest struct {
	Symbols []string `json:"symbols" validate:"required,min=1"`
}

type PriceUpdate struct {
	Symbol    string  `json:"symbol"`
	Price     float64 `json:"price"`
	Volume    float64 `json:"volume"`
	Timestamp int64   `json:"timestamp"`
}
