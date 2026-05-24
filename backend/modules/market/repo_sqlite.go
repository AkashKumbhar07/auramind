package market

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/AkashKumbhar07/auramind/backend/framework/errors"
)

type SQLiteRepository struct {
	db *sql.DB
}

func NewSQLiteRepository(db *sql.DB) *SQLiteRepository {
	return &SQLiteRepository{db: db}
}

func (r *SQLiteRepository) GetPairs(ctx context.Context, exchange string) ([]MarketPair, error) {
	query := `SELECT id, symbol, base_asset, quote_asset, exchange, active FROM market_pairs`
	args := []any{}
	if exchange != "" {
		query += ` WHERE exchange = ?`
		args = append(args, exchange)
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, errors.Wrap(errors.KindInternal, "get pairs", err)
	}
	defer rows.Close()

	var pairs []MarketPair
	for rows.Next() {
		var p MarketPair
		if err := rows.Scan(&p.ID, &p.Symbol, &p.BaseAsset, &p.QuoteAsset, &p.Exchange, &p.Active); err != nil {
			return nil, errors.Wrap(errors.KindInternal, "scan pair", err)
		}
		pairs = append(pairs, p)
	}
	return pairs, nil
}

func (r *SQLiteRepository) GetCandles(ctx context.Context, symbol, exchange, interval string, limit int) ([]Candle, error) {
	query := `SELECT symbol, exchange, interval, open, high, low, close, volume, timestamp
	           FROM candles WHERE symbol = ? AND exchange = ? AND interval = ?
	           ORDER BY timestamp DESC LIMIT ?`

	rows, err := r.db.QueryContext(ctx, query, symbol, exchange, interval, limit)
	if err != nil {
		return nil, errors.Wrap(errors.KindInternal, "get candles", err)
	}
	defer rows.Close()

	var candles []Candle
	for rows.Next() {
		var c Candle
		if err := rows.Scan(&c.Symbol, &c.Exchange, &c.Interval, &c.Open, &c.High, &c.Low, &c.Close, &c.Volume, &c.Timestamp); err != nil {
			return nil, errors.Wrap(errors.KindInternal, "scan candle", err)
		}
		candles = append(candles, c)
	}
	return candles, nil
}

func (r *SQLiteRepository) SaveCandle(ctx context.Context, candle *Candle) error {
	query := `INSERT OR REPLACE INTO candles (symbol, exchange, interval, open, high, low, close, volume, timestamp)
	           VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`

	_, err := r.db.ExecContext(ctx, query,
		candle.Symbol, candle.Exchange, candle.Interval,
		candle.Open, candle.High, candle.Low, candle.Close, candle.Volume, candle.Timestamp,
	)
	if err != nil {
		return errors.Wrap(errors.KindInternal, "save candle", err)
	}
	return nil
}

func (r *SQLiteRepository) BulkSaveCandles(ctx context.Context, candles []*Candle) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return errors.Wrap(errors.KindInternal, "begin tx", err)
	}
	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx, `INSERT OR REPLACE INTO candles (symbol, exchange, interval, open, high, low, close, volume, timestamp) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`)
	if err != nil {
		return fmt.Errorf("prepare: %w", err)
	}
	defer stmt.Close()

	for _, c := range candles {
		_, err := stmt.ExecContext(ctx, c.Symbol, c.Exchange, c.Interval, c.Open, c.High, c.Low, c.Close, c.Volume, c.Timestamp)
		if err != nil {
			return errors.Wrap(errors.KindInternal, "insert candle", err)
		}
	}

	return tx.Commit()
}
