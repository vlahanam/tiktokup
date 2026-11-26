package middleware

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"strings"

	"tiktokup/backend/pkg/response"
)

const RequestIDHeader = "X-Request-ID"

// RequestID ensures every incoming request has a correlation ID propagated via header and context.
func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := strings.TrimSpace(r.Header.Get(RequestIDHeader))
		if id == "" || len(id) > 128 {
			id = generateRequestID()
		}

		w.Header().Set(RequestIDHeader, id)
		ctx := response.WithRequestID(r.Context(), id)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func generateRequestID() string {
	var b [16]byte
	if _, err := rand.Read(b[:]); err != nil {
		return "req-fallback"
	}
	return hex.EncodeToString(b[:])
}
