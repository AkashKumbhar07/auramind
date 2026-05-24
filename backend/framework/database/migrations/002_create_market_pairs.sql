-- +goose Up
CREATE TABLE market_pairs (
    id          TEXT PRIMARY KEY,
    symbol      TEXT NOT NULL,
    base_asset  TEXT NOT NULL,
    quote_asset TEXT NOT NULL,
    exchange    TEXT NOT NULL,
    active      INTEGER NOT NULL DEFAULT 1,
    created_at  INTEGER NOT NULL DEFAULT (strftime('%s', 'now')),
    updated_at  INTEGER NOT NULL DEFAULT (strftime('%s', 'now')),
    UNIQUE(symbol, exchange)
);

CREATE INDEX idx_market_pairs_symbol ON market_pairs(symbol);
CREATE INDEX idx_market_pairs_exchange ON market_pairs(exchange);

-- +goose Down
DROP TABLE IF EXISTS market_pairs;
