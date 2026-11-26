package config

import (
	"log/slog"
	"testing"
	"time"
)

func TestLoadDefaults(t *testing.T) {
	t.Setenv("APP_NAME", "")
	t.Setenv("APP_ENV", "")
	t.Setenv("PORT", "")
	t.Setenv("LOG_LEVEL", "")
	t.Setenv("CORS_ALLOWED_ORIGINS", "")
	t.Setenv("SERVER_READ_TIMEOUT_SECONDS", "")
	t.Setenv("DATABASE_URL", "")
	t.Setenv("DB_HOST", "")
	t.Setenv("DB_PORT", "")
	t.Setenv("DB_USER", "")
	t.Setenv("DB_PASSWORD", "")
	t.Setenv("DB_NAME", "")
	t.Setenv("DB_SSLMODE", "")

	cfg := Load()

	if cfg.AppName != "tiktokup-backend" {
		t.Fatalf("expected default AppName tiktokup-backend, got %q", cfg.AppName)
	}
	if cfg.Addr() != ":8080" {
		t.Fatalf("expected default Addr :8080, got %q", cfg.Addr())
	}
	if cfg.SlogLevel() != slog.LevelInfo {
		t.Fatalf("expected default SlogLevel info, got %v", cfg.SlogLevel())
	}
	if cfg.ReadTimeout != 15*time.Second {
		t.Fatalf("expected ReadTimeout 15s, got %v", cfg.ReadTimeout)
	}
	expectedDSN := "postgres://tiktokup:tiktokup_secret@postgres:5432/tiktokup?sslmode=disable"
	if cfg.Database.DSN() != expectedDSN {
		t.Fatalf("expected default DSN %q, got %q", expectedDSN, cfg.Database.DSN())
	}
	if !cfg.Database.AutoMigrate {
		t.Fatalf("expected default AutoMigrate true")
	}
}

func TestLoadOverrides(t *testing.T) {
	t.Setenv("APP_NAME", "tiktokup-api")
	t.Setenv("APP_ENV", "production")
	t.Setenv("PORT", "9090")
	t.Setenv("LOG_LEVEL", "debug")
	t.Setenv("CORS_ALLOWED_ORIGINS", "https://example.com, https://app.example.com")
	t.Setenv("SERVER_READ_TIMEOUT_SECONDS", "25")
	t.Setenv("DATABASE_URL", "postgres://prod:secret@db.internal:5432/prod_db?sslmode=require")
	t.Setenv("DB_MAX_CONNS", "25")
	t.Setenv("DB_AUTO_MIGRATE", "false")

	cfg := Load()

	if cfg.AppName != "tiktokup-api" {
		t.Fatalf("unexpected AppName: %q", cfg.AppName)
	}
	if cfg.AppEnv != "production" {
		t.Fatalf("unexpected AppEnv: %q", cfg.AppEnv)
	}
	if cfg.Addr() != ":9090" {
		t.Fatalf("unexpected Addr: %q", cfg.Addr())
	}
	if cfg.SlogLevel() != slog.LevelDebug {
		t.Fatalf("expected debug level, got %v", cfg.SlogLevel())
	}
	if len(cfg.CORSAllowedOrigins) != 2 || cfg.CORSAllowedOrigins[1] != "https://app.example.com" {
		t.Fatalf("unexpected CORSAllowedOrigins: %#v", cfg.CORSAllowedOrigins)
	}
	if cfg.ReadTimeout != 25*time.Second {
		t.Fatalf("expected ReadTimeout 25s, got %v", cfg.ReadTimeout)
	}
	if cfg.Database.DSN() != "postgres://prod:secret@db.internal:5432/prod_db?sslmode=require" {
		t.Fatalf("unexpected DATABASE_URL override: %q", cfg.Database.DSN())
	}
	if cfg.Database.MaxConns != 25 || cfg.Database.AutoMigrate {
		t.Fatalf("unexpected DB overrides: max=%d auto=%v", cfg.Database.MaxConns, cfg.Database.AutoMigrate)
	}
}
