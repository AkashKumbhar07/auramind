package config

import (
	"path/filepath"
	"testing"
)

func projectRoot() string {
	// Walk up from the test file to find project root
	return filepath.Join("..", "..", "..")
}

func TestLoad_Development(t *testing.T) {
	cfg, err := Load(filepath.Join(projectRoot(), "configs", "environments", "development.yaml"))
	if err != nil {
		t.Fatalf("failed to load config: %v", err)
	}

	if cfg.App.Name != "AuraMind AI" {
		t.Errorf("expected AuraMind AI, got %s", cfg.App.Name)
	}

	if cfg.DB.Postgres.Host != "localhost" {
		t.Errorf("expected localhost, got %s", cfg.DB.Postgres.Host)
	}

	if cfg.Msg.Provider != "nats" {
		t.Errorf("expected nats, got %s", cfg.Msg.Provider)
	}
}

func TestPostgresDSN(t *testing.T) {
	p := PostgresConfig{
		Host: "localhost", Port: 5432, User: "u",
		Password: "p", Database: "d", SSLMode: "disable",
	}

	dsn := p.DSN()
	if dsn != "host=localhost port=5432 user=u password=p dbname=d sslmode=disable" {
		t.Errorf("unexpected DSN: %s", dsn)
	}
}
