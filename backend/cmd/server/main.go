package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"tiktokup/backend/internal/config"
	"tiktokup/backend/internal/database"
	"tiktokup/backend/internal/pythonclient"
	"tiktokup/backend/internal/router"
	"tiktokup/backend/internal/server"
	"tiktokup/backend/migrations"
)

func main() {
	cfg := config.Load()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: cfg.SlogLevel(),
	}))
	slog.SetDefault(logger)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	db, err := database.Connect(ctx, cfg.Database, logger)
	if err != nil {
		logger.Error("postgres_connect_fatal", slog.String("error", err.Error()))
		os.Exit(1)
	}
	defer db.Close()

	if cfg.Database.AutoMigrate {
		applied, err := database.Migrate(ctx, db.Pool, migrations.Files, logger)
		if err != nil {
			logger.Error("postgres_migrate_fatal", slog.String("error", err.Error()))
			os.Exit(1)
		}
		logger.Info("postgres_migrations_ready", slog.Int("applied", applied))
	}

	pyClient := pythonclient.New(cfg.PythonServiceURL)
	handler := router.New(cfg, logger, db, pyClient)
	srv := server.New(cfg, handler, logger)

	if err := srv.Run(ctx); err != nil {
		logger.Error("http_server_fatal", slog.String("error", err.Error()))
		os.Exit(1)
	}
}
