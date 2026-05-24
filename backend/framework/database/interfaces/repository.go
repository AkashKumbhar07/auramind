package interfaces

import (
	"context"
)

type Repository[T any] interface {
	Create(ctx context.Context, item *T) error
	GetByID(ctx context.Context, id string) (*T, error)
	Update(ctx context.Context, item *T) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, filter map[string]any) ([]*T, error)
}

type UserRepository interface {
	Repository[User]
	GetByEmail(ctx context.Context, email string) (*User, error)
	GetByUsername(ctx context.Context, username string) (*User, error)
}

type MarketRepository interface {
	Repository[MarketPair]
	GetBySymbol(ctx context.Context, symbol string) (*MarketPair, error)
	GetActive(ctx context.Context) ([]*MarketPair, error)
}

type CandleRepository interface {
	BulkCreate(ctx context.Context, candles []*Candle) error
	GetBySymbol(ctx context.Context, symbol string, limit int) ([]*Candle, error)
	GetRange(ctx context.Context, symbol string, from, to int64) ([]*Candle, error)
}

type AlertRepository interface {
	Repository[Alert]
	GetByUserID(ctx context.Context, userID string) ([]*Alert, error)
	GetActive(ctx context.Context) ([]*Alert, error)
}

type StrategyRepository interface {
	Repository[Strategy]
	GetByUserID(ctx context.Context, userID string) ([]*Strategy, error)
	GetPublic(ctx context.Context) ([]*Strategy, error)
}

type TradeRepository interface {
	Repository[Trade]
	GetByUserID(ctx context.Context, userID string) ([]*Trade, error)
	GetBySymbol(ctx context.Context, symbol string) ([]*Trade, error)
}

type PortfolioRepository interface {
	Repository[Portfolio]
	GetByUserID(ctx context.Context, userID string) (*Portfolio, error)
}
