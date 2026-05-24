package market

import "context"

type Repository interface {
	GetPairs(ctx context.Context, exchange string) ([]MarketPair, error)
	GetCandles(ctx context.Context, symbol, exchange, interval string, limit int) ([]Candle, error)
	SaveCandle(ctx context.Context, candle *Candle) error
	BulkSaveCandles(ctx context.Context, candles []*Candle) error
}
