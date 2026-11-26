package router

import (
	"log/slog"
	"net/http"

	"tiktokup/backend/internal/config"
	"tiktokup/backend/internal/handler"
	"tiktokup/backend/internal/middleware"
	"tiktokup/backend/pkg/response"
)

// New builds the HTTP handler tree with versioned routes and standard middleware.
func New(cfg *config.Config, logger *slog.Logger, db handler.DatabaseInspector) http.Handler {
	mux := http.NewServeMux()
	sys := handler.NewSystemHandler(cfg, db)

	// Container & orchestrator probes
	mux.HandleFunc("GET /health", sys.Health)

	// API v1 system endpoints
	mux.HandleFunc("GET /api/v1/health", sys.Health)
	mux.HandleFunc("GET /api/v1/ready", sys.Ready)
	mux.HandleFunc("GET /api/v1/info", sys.Info)

	// Structured JSON fallback for unmatched routes
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		response.Error(w, r, http.StatusNotFound, "NOT_FOUND", "endpoint not found")
	})

	var h http.Handler = mux
	h = middleware.CORS(cfg.CORSAllowedOrigins)(h)
	h = middleware.Logger(logger)(h)
	h = middleware.RequestID(h)
	h = middleware.Recover(logger)(h)
	return h
}
