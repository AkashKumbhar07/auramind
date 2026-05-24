package market

const (
	EventCandleCreated  = "market.candle.created"
	EventTradeExecuted  = "market.trade.executed"
	EventPriceUpdate    = "market.price.updated"
	EventAnalysisReady  = "market.analysis.ready"
)

type CandleCreatedEvent struct {
	Symbol   string  `json:"symbol"`
	Exchange string  `json:"exchange"`
	Interval string  `json:"interval"`
	Close    float64 `json:"close"`
	Volume   float64 `json:"volume"`
}

type PriceUpdateEvent struct {
	Symbol    string  `json:"symbol"`
	Exchange  string  `json:"exchange"`
	Price     float64 `json:"price"`
	Volume    float64 `json:"volume"`
	Timestamp int64   `json:"timestamp"`
}
