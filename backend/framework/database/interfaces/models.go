package interfaces

type User struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	Username  string `json:"username"`
	Password  string `json:"-"`
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
}

type MarketPair struct {
	ID            string `json:"id"`
	Symbol        string `json:"symbol"`
	BaseAsset     string `json:"base_asset"`
	QuoteAsset    string `json:"quote_asset"`
	Exchange      string `json:"exchange"`
	Active        bool   `json:"active"`
	CreatedAt     int64  `json:"created_at"`
	UpdatedAt     int64  `json:"updated_at"`
}

type Candle struct {
	ID        string  `json:"id"`
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

type Alert struct {
	ID        string `json:"id"`
	UserID    string `json:"user_id"`
	Symbol    string `json:"symbol"`
	Condition string `json:"condition"`
	Value     float64 `json:"value"`
	Active    bool   `json:"active"`
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
}

type Strategy struct {
	ID        string  `json:"id"`
	UserID    string  `json:"user_id"`
	Name      string  `json:"name"`
	Config    string  `json:"config"`
	Public    bool    `json:"public"`
	PnL       float64 `json:"pnl"`
	CreatedAt int64   `json:"created_at"`
	UpdatedAt int64   `json:"updated_at"`
}

type Trade struct {
	ID        string  `json:"id"`
	UserID    string  `json:"user_id"`
	Symbol    string  `json:"symbol"`
	Side      string  `json:"side"`
	Amount    float64 `json:"amount"`
	Price     float64 `json:"price"`
	PnL       float64 `json:"pnl"`
	Status    string  `json:"status"`
	CreatedAt int64   `json:"created_at"`
}

type Portfolio struct {
	ID         string             `json:"id"`
	UserID     string             `json:"user_id"`
	Positions  []PortfolioPosition `json:"positions"`
	TotalValue float64            `json:"total_value"`
	CreatedAt  int64              `json:"created_at"`
	UpdatedAt  int64              `json:"updated_at"`
}

type PortfolioPosition struct {
	Symbol   string  `json:"symbol"`
	Amount   float64 `json:"amount"`
	Entry    float64 `json:"entry"`
	Current  float64 `json:"current"`
	UnrealPnL float64 `json:"unrealized_pnl"`
}
