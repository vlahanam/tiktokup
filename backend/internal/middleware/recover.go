package middleware

import (
	"log/slog"
	"net/http"
	"runtime/debug"

	"tiktokup/backend/pkg/response"
)

// Recover catches panics in downstream handlers, logs the stack trace, and responds with HTTP 500 JSON.
func Recover(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rec := recover(); rec != nil {
					logger.ErrorContext(r.Context(), "http_panic_recovered",
						slog.Any("panic", rec),
						slog.String("request_id", response.RequestIDFromContext(r.Context())),
						slog.String("method", r.Method),
						slog.String("path", r.URL.Path),
						slog.String("stack", string(debug.Stack())),
					)
					response.Error(w, r, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "an unexpected internal error occurred")
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}
