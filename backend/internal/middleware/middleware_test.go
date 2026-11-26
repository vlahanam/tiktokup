package middleware

import (
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"tiktokup/backend/pkg/response"
)

func TestRequestIDMiddleware(t *testing.T) {
	var seenID string
	h := RequestID(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		seenID = response.RequestIDFromContext(r.Context())
		w.WriteHeader(http.StatusOK)
	}))

	// Auto-generated ID
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if seenID == "" || rec.Header().Get(RequestIDHeader) != seenID {
		t.Fatalf("expected generated request id in context and header, got ctx=%q hdr=%q", seenID, rec.Header().Get(RequestIDHeader))
	}

	// Preserved incoming ID
	req2 := httptest.NewRequest(http.MethodGet, "/health", nil)
	req2.Header.Set(RequestIDHeader, "trace-abc-999")
	rec2 := httptest.NewRecorder()
	h.ServeHTTP(rec2, req2)

	if seenID != "trace-abc-999" || rec2.Header().Get(RequestIDHeader) != "trace-abc-999" {
		t.Fatalf("expected preserved request id trace-abc-999, got %q", seenID)
	}
}

func TestCORSMiddleware(t *testing.T) {
	h := CORS([]string{"http://localhost:8000"})(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	// Preflight allowed origin
	req := httptest.NewRequest(http.MethodOptions, "/api/v1/health", nil)
	req.Header.Set("Origin", "http://localhost:8000")
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("expected 204 for preflight, got %d", rec.Code)
	}
	if got := rec.Header().Get("Access-Control-Allow-Origin"); got != "http://localhost:8000" {
		t.Fatalf("expected allowed origin header, got %q", got)
	}

	// Disallowed origin
	req2 := httptest.NewRequest(http.MethodGet, "/api/v1/health", nil)
	req2.Header.Set("Origin", "https://evil.example")
	rec2 := httptest.NewRecorder()
	h.ServeHTTP(rec2, req2)

	if got := rec2.Header().Get("Access-Control-Allow-Origin"); got != "" {
		t.Fatalf("expected empty allow-origin for disallowed origin, got %q", got)
	}
}

func TestRecoverMiddleware(t *testing.T) {
	logger := slog.New(slog.NewJSONHandler(io.Discard, nil))
	h := Recover(logger)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("boom")
	}))

	req := httptest.NewRequest(http.MethodGet, "/panic", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500 after panic, got %d", rec.Code)
	}
}
