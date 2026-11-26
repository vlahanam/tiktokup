package config

import (
	"fmt"
	"log/slog"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

// DatabaseConfig holds PostgreSQL connection and pool settings.
type DatabaseConfig struct {
	URL         string
	Host        string
	Port        string
	User        string
	Password    string
	Name        string
	SSLMode     string
	MaxConns    int32
	MinConns    int32
	AutoMigrate bool
}

// DSN builds a PostgreSQL connection URI. If DATABASE_URL is explicitly set, it takes precedence.
func (d DatabaseConfig) DSN() string {
	if d.URL != "" {
		return d.URL
	}
	u := &url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(d.User, d.Password),
		Host:   fmt.Sprintf("%s:%s", d.Host, d.Port),
		Path:   d.Name,
	}
	q := u.Query()
	if d.SSLMode != "" {
		q.Set("sslmode", d.SSLMode)
	}
	u.RawQuery = q.Encode()
	return u.String()
}

// Config holds all runtime configuration loaded from environment variables.
type Config struct {
	AppName            string
	AppVersion         string
	AppEnv             string
	Port               string
	LogLevel           string
	CORSAllowedOrigins []string
	ReadTimeout        time.Duration
	WriteTimeout       time.Duration
	IdleTimeout        time.Duration
	ShutdownTimeout    time.Duration
	Database           DatabaseConfig
}

// Load parses environment variables and returns a validated Config with sensible defaults.
func Load() *Config {
	return &Config{
		AppName:            getEnv("APP_NAME", "tiktokup-backend"),
		AppVersion:         getEnv("APP_VERSION", "0.1.0"),
		AppEnv:             getEnv("APP_ENV", "development"),
		Port:               getEnv("PORT", "8080"),
		LogLevel:           strings.ToLower(getEnv("LOG_LEVEL", "info")),
		CORSAllowedOrigins: parseCSV(getEnv("CORS_ALLOWED_ORIGINS", "http://localhost:8000,http://localhost:3000")),
		ReadTimeout:        getEnvDurationSeconds("SERVER_READ_TIMEOUT_SECONDS", 15),
		WriteTimeout:       getEnvDurationSeconds("SERVER_WRITE_TIMEOUT_SECONDS", 30),
		IdleTimeout:        getEnvDurationSeconds("SERVER_IDLE_TIMEOUT_SECONDS", 60),
		ShutdownTimeout:    getEnvDurationSeconds("SERVER_SHUTDOWN_TIMEOUT_SECONDS", 10),
		Database: DatabaseConfig{
			URL:         getEnv("DATABASE_URL", ""),
			Host:        getEnv("DB_HOST", "postgres"),
			Port:        getEnv("DB_PORT", "5432"),
			User:        getEnv("DB_USER", "tiktokup"),
			Password:    getEnv("DB_PASSWORD", "tiktokup_secret"),
			Name:        getEnv("DB_NAME", "tiktokup"),
			SSLMode:     getEnv("DB_SSLMODE", "disable"),
			MaxConns:    int32(getEnvInt("DB_MAX_CONNS", 10)),
			MinConns:    int32(getEnvInt("DB_MIN_CONNS", 2)),
			AutoMigrate: getEnvBool("DB_AUTO_MIGRATE", true),
		},
	}
}

// Addr returns the listen address formatted for net/http.
func (c *Config) Addr() string {
	if strings.HasPrefix(c.Port, ":") {
		return c.Port
	}
	return ":" + c.Port
}

// SlogLevel converts the string log level into a slog.Level.
func (c *Config) SlogLevel() slog.Level {
	switch c.LogLevel {
	case "debug":
		return slog.LevelDebug
	case "warn", "warning":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

func getEnv(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok {
		trimmed := strings.TrimSpace(val)
		if trimmed != "" {
			return trimmed
		}
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	raw := getEnv(key, "")
	if raw == "" {
		return fallback
	}
	n, err := strconv.Atoi(raw)
	if err != nil || n <= 0 {
		return fallback
	}
	return n
}

func getEnvBool(key string, fallback bool) bool {
	raw := getEnv(key, "")
	if raw == "" {
		return fallback
	}
	b, err := strconv.ParseBool(raw)
	if err != nil {
		return fallback
	}
	return b
}

func getEnvDurationSeconds(key string, fallbackSeconds int) time.Duration {
	return time.Duration(getEnvInt(key, fallbackSeconds)) * time.Second
}

func parseCSV(raw string) []string {
	parts := strings.Split(raw, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if trimmed := strings.TrimSpace(p); trimmed != "" {
			out = append(out, trimmed)
		}
	}
	return out
}
