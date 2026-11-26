package router_test

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"tiktokup/backend/internal/config"
	"tiktokup/backend/internal/handler"
	"tiktokup/backend/internal/router"
)

type stubDB struct {
	pingErr error
	version string
}

func (s stubDB) Ping(context.Context) error { return s.pingErr }
func (s stubDB) ServerVersion(context.Context) (string, error) {
	if s.pingErr != nil {
		return "", s.pingErr
	}
	return s.version, nil
}

func newTestRouter(db handler.DatabaseInspector) http.Handler {
	cfg := config.Load()
	logger := slog.New(slog.NewJSONHandler(io.Discard, nil))
	return router.New(cfg, logger, db)
}

func TestSystemEndpoints(t *testing.T) {
	h := newTestRouter(stubDB{version: "PostgreSQL 18.6"})

	tests := []struct {
		name           string
		method         string
		path           string
		expectedStatus int
		expectSuccess  bool
	}{
		{"root health", http.MethodGet, "/health", http.StatusOK, true},
		{"api v1 health", http.MethodGet, "/api/v1/health", http.StatusOK, true},
		{"api v1 ready", http.MethodGet, "/api/v1/ready", http.StatusOK, true},
		{"api v1 info", http.MethodGet, "/api/v1/info", http.StatusOK, true},
		{"unknown api route", http.MethodGet, "/api/v1/does-not-exist", http.StatusNotFound, false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(tc.method, tc.path, nil)
			rec := httptest.NewRecorder()
			h.ServeHTTP(rec, req)

			if rec.Code != tc.expectedStatus {
				t.Fatalf("expected status %d, got %d (body=%s)", tc.expectedStatus, rec.Code, rec.Body.String())
			}

			var body struct {
				Success bool            `json:"success"`
				Data    json.RawMessage `json:"data"`
			}
			if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
				t.Fatalf("failed to decode json: %v", err)
			}
			if body.Success != tc.expectSuccess {
				t.Fatalf("expected success=%v, got %v", tc.expectSuccess, body.Success)
			}
			if rec.Header().Get("X-Request-ID") == "" {
				t.Fatalf("expected X-Request-ID header to be set")
			}
		})
	}
}

func TestInfoReturnsGoAndPostgresMetadata(t *testing.T) {
	h := newTestRouter(stubDB{version: "PostgreSQL 18.6"})
	req := httptest.NewRequest(http.MethodGet, "/api/v1/info", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	body := rec.Body.String()
	if !strings.Contains(body, "go1.26") {
		t.Fatalf("expected response to report go1.26 runtime, got: %s", body)
	}
	if !strings.Contains(body, "PostgreSQL 18.6") {
		t.Fatalf("expected response to report PostgreSQL 18.6, got: %s", body)
	}
}

func TestReadyReturns503WhenDatabaseDown(t *testing.T) {
	h := newTestRouter(stubDB{pingErr: errors.New("connection refused")})
	req := httptest.NewRequest(http.MethodGet, "/api/v1/ready", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusServiceUnavailable {
		t.Fatalf("expected 503 when DB is down, got %d (body=%s)", rec.Code, rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), `"postgres":"error"`) {
		t.Fatalf("expected postgres error check in payload, got: %s", rec.Body.String())
	}
}
