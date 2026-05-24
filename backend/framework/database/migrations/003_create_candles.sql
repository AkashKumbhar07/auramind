-- +goose Up
CREATE TABLE candles (
    id         TEXT PRIMARY KEY,
    symbol     TEXT NOT NULL,
    exchange   TEXT NOT NULL,
    interval   TEXT NOT NULL,
    open       REAL NOT NULL,
    high       REAL NOT NULL,
    low        REAL NOT NULL,
    close      REAL NOT NULL,
    volume     REAL NOT NULL,
    timestamp  INTEGER NOT NULL,
    UNIQUE(symbol, exchange, interval, timestamp)
);

CREATE INDEX idx_candles_symbol ON candles(symbol);
CREATE INDEX idx_candles_timestamp ON candles(timestamp);
CREATE INDEX idx_candles_symbol_ts ON candles(symbol, timestamp);

-- +goose Down
DROP TABLE IF EXISTS candles;
