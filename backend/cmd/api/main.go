package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/AkashKumbhar07/auramind/backend/framework/config"
	"github.com/AkashKumbhar07/auramind/backend/framework/database"
	grpcserver "github.com/AkashKumbhar07/auramind/backend/framework/grpc/server"
	"github.com/AkashKumbhar07/auramind/backend/framework/logging"
	marketgrpc "github.com/AkashKumbhar07/auramind/backend/framework/grpc/generated/market"
	marketmodule "github.com/AkashKumbhar07/auramind/backend/modules/market"
	"go.uber.org/zap"
	_ "modernc.org/sqlite"
)

func main() {
	// Config
	configPath := filepath.Join("configs", "environments", "development.yaml")
	cfg, err := config.Load(configPath)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// Logger
	logger := logging.New(cfg.Log.Level, cfg.Log.Format)
	defer logger.Sync()

	logger.Info("starting API gateway",
		zap.String("app", cfg.App.Name),
		zap.String("env", cfg.App.Env),
		zap.Int("port", cfg.GRPC.Port),
	)

	// Database
	dbProvider := database.ProviderSQLite
	db, err := database.Open(cfg.DB, dbProvider)
	if err != nil {
		logger.Fatal("failed to connect to database", zap.Error(err))
	}
	defer db.Close()
	logger.Info("database connected", zap.String("provider", string(dbProvider)))

	// Repository
	repo := marketmodule.NewSQLiteRepository(db)

	// Market service
	marketSvc := marketmodule.NewService(repo, logger)

	// gRPC server
	grpcServer := grpcserver.New(cfg.GRPC.Port, logger)

	// Register services
	marketGRPC := marketmodule.NewGRPCServer(marketSvc)
	marketgrpc.RegisterMarketServiceServer(grpcServer.Server, marketGRPC)
	logger.Info("registered market service")

	// Start gRPC
	go func() {
		if err := grpcServer.Start(); err != nil {
			logger.Fatal("gRPC server failed", zap.Error(err))
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("shutting down...")
	grpcServer.Shutdown()
	logger.Info("server stopped")
	_ = context.Background()
}
