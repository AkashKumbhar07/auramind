package database

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	_ "modernc.org/sqlite"

	"github.com/AkashKumbhar07/auramind/backend/framework/config"
)

func OpenPostgres(cfg config.PostgresConfig) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.DSN())
	if err != nil {
		return nil, fmt.Errorf("open postgres: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("ping postgres: %w", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)

	return db, nil
}

func OpenSQLite(cfg config.SQLiteConfig) (*sql.DB, error) {
	db, err := sql.Open("sqlite", cfg.Path)
	if err != nil {
		return nil, fmt.Errorf("open sqlite: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("ping sqlite: %w", err)
	}

	db.SetMaxOpenConns(1)

	return db, nil
}

type Provider string

const (
	ProviderPostgres Provider = "postgres"
	ProviderSQLite   Provider = "sqlite"
)

func Open(cfg config.DBConfig, provider Provider) (*sql.DB, error) {
	switch provider {
	case ProviderPostgres:
		return OpenPostgres(cfg.Postgres)
	case ProviderSQLite:
		return OpenSQLite(cfg.SQLite)
	default:
		return nil, fmt.Errorf("unknown database provider: %s", provider)
	}
}
