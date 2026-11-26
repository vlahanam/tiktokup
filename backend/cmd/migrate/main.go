package main

import (
	"context"
	"log/slog"
	"os"
	"time"

	"tiktokup/backend/internal/config"
	"tiktokup/backend/internal/database"
	"tiktokup/backend/migrations"
)

func main() {
	cfg := config.Load()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: cfg.SlogLevel(),
	}))
	slog.SetDefault(logger)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	db, err := database.Connect(ctx, cfg.Database, logger)
	if err != nil {
		logger.Error("migrate_connect_failed", slog.String("error", err.Error()))
		os.Exit(1)
	}
	defer db.Close()

	applied, err := database.Migrate(ctx, db.Pool, migrations.Files, logger)
	if err != nil {
		logger.Error("migrate_failed", slog.String("error", err.Error()))
		os.Exit(1)
	}

	logger.Info("migrate_complete", slog.Int("applied", applied))
}
