package database

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"tiktokup/backend/internal/config"
)

// DB wraps a pgxpool.Pool and exposes health and diagnostic helpers.
type DB struct {
	Pool *pgxpool.Pool
}

// Connect creates a configured PostgreSQL connection pool and verifies connectivity with retries.
func Connect(ctx context.Context, cfg config.DatabaseConfig, logger *slog.Logger) (*DB, error) {
	poolCfg, err := pgxpool.ParseConfig(cfg.DSN())
	if err != nil {
		return nil, fmt.Errorf("parse database config: %w", err)
	}

	if cfg.MaxConns > 0 {
		poolCfg.MaxConns = cfg.MaxConns
	}
	if cfg.MinConns > 0 && cfg.MinConns <= poolCfg.MaxConns {
		poolCfg.MinConns = cfg.MinConns
	}
	poolCfg.MaxConnLifetime = 30 * time.Minute
	poolCfg.MaxConnIdleTime = 5 * time.Minute
	poolCfg.HealthCheckPeriod = 30 * time.Second

	pool, err := pgxpool.NewWithConfig(ctx, poolCfg)
	if err != nil {
		return nil, fmt.Errorf("create pgx pool: %w", err)
	}

	var lastErr error
	const maxAttempts = 5
	for attempt := 1; attempt <= maxAttempts; attempt++ {
		pingCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
		lastErr = pool.Ping(pingCtx)
		cancel()
		if lastErr == nil {
			logger.Info("postgres_connected",
				slog.String("host", cfg.Host),
				slog.String("port", cfg.Port),
				slog.String("database", cfg.Name),
				slog.Int("max_conns", int(poolCfg.MaxConns)),
			)
			return &DB{Pool: pool}, nil
		}

		logger.Warn("postgres_ping_retry",
			slog.Int("attempt", attempt),
			slog.Int("max_attempts", maxAttempts),
			slog.String("error", lastErr.Error()),
		)

		select {
		case <-ctx.Done():
			pool.Close()
			return nil, ctx.Err()
		case <-time.After(time.Duration(attempt) * time.Second):
		}
	}

	pool.Close()
	return nil, fmt.Errorf("connect to postgres after %d attempts: %w", maxAttempts, lastErr)
}

// Ping verifies live connectivity to PostgreSQL.
func (d *DB) Ping(ctx context.Context) error {
	if d == nil || d.Pool == nil {
		return fmt.Errorf("database pool is not initialized")
	}
	return d.Pool.Ping(ctx)
}

// ServerVersion queries PostgreSQL for its concise version string (e.g. "PostgreSQL 18.6").
func (d *DB) ServerVersion(ctx context.Context) (string, error) {
	if d == nil || d.Pool == nil {
		return "", fmt.Errorf("database pool is not initialized")
	}
	var raw string
	if err := d.Pool.QueryRow(ctx, "SELECT version()").Scan(&raw); err != nil {
		return "", err
	}
	// Extract the concise leading tokens, e.g. "PostgreSQL 18.6"
	fields := strings.Fields(raw)
	if len(fields) >= 2 {
		return fields[0] + " " + fields[1], nil
	}
	return raw, nil
}

// Close releases all pooled database connections.
func (d *DB) Close() {
	if d != nil && d.Pool != nil {
		d.Pool.Close()
	}
}
