package handler

import (
	"context"
	"net/http"
	"runtime"
	"time"

	"tiktokup/backend/internal/config"
	"tiktokup/backend/pkg/response"
)

// DatabaseInspector abstracts live database health and metadata probes for handlers and tests.
type DatabaseInspector interface {
	Ping(ctx context.Context) error
	ServerVersion(ctx context.Context) (string, error)
}

// SystemHandler handles health, readiness, and service metadata endpoints.
type SystemHandler struct {
	cfg       *config.Config
	db        DatabaseInspector
	startedAt time.Time
}

// NewSystemHandler constructs a SystemHandler anchored to startup time.
func NewSystemHandler(cfg *config.Config, db DatabaseInspector) *SystemHandler {
	return &SystemHandler{
		cfg:       cfg,
		db:        db,
		startedAt: time.Now().UTC(),
	}
}

// HealthResponse represents the liveness check payload.
type HealthResponse struct {
	Status        string `json:"status"`
	Service       string `json:"service"`
	Environment   string `json:"environment"`
	UptimeSeconds int64  `json:"uptime_seconds"`
}

// ReadyResponse represents the readiness probe payload.
type ReadyResponse struct {
	Status string            `json:"status"`
	Checks map[string]string `json:"checks"`
}

// DatabaseInfo exposes PostgreSQL engine status and version metadata.
type DatabaseInfo struct {
	Engine    string `json:"engine"`
	Connected bool   `json:"connected"`
	Version   string `json:"version,omitempty"`
}

// InfoResponse exposes runtime metadata for diagnostics and status UIs.
type InfoResponse struct {
	Service       string        `json:"service"`
	Version       string        `json:"version"`
	Environment   string        `json:"environment"`
	GoVersion     string        `json:"go_version"`
	OS            string        `json:"os"`
	Arch          string        `json:"arch"`
	NumCPU        int           `json:"num_cpu"`
	StartedAt     string        `json:"started_at"`
	UptimeSeconds int64         `json:"uptime_seconds"`
	Database      *DatabaseInfo `json:"database,omitempty"`
}

// Health responds with service liveness state.
func (h *SystemHandler) Health(w http.ResponseWriter, r *http.Request) {
	response.OK(w, r, HealthResponse{
		Status:        "ok",
		Service:       h.cfg.AppName,
		Environment:   h.cfg.AppEnv,
		UptimeSeconds: int64(time.Since(h.startedAt).Seconds()),
	})
}

// Ready responds with readiness probe status including live PostgreSQL connectivity.
func (h *SystemHandler) Ready(w http.ResponseWriter, r *http.Request) {
	checks := map[string]string{
		"http_server": "ok",
	}

	if h.db != nil {
		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()

		if err := h.db.Ping(ctx); err != nil {
			checks["postgres"] = "error"
			response.WriteJSON(w, http.StatusServiceUnavailable, response.Envelope{
				Success: false,
				Data: ReadyResponse{
					Status: "not_ready",
					Checks: checks,
				},
				Error: &response.ErrorInfo{
					Code:    "SERVICE_UNAVAILABLE",
					Message: "database connectivity check failed",
				},
				Meta: response.Meta{
					RequestID: response.RequestIDFromContext(r.Context()),
					Timestamp: time.Now().UTC().Format(time.RFC3339),
				},
			})
			return
		}
		checks["postgres"] = "ok"
	}

	response.OK(w, r, ReadyResponse{
		Status: "ready",
		Checks: checks,
	})
}

// Info responds with Go runtime, service, and PostgreSQL metadata.
func (h *SystemHandler) Info(w http.ResponseWriter, r *http.Request) {
	var dbInfo *DatabaseInfo
	if h.db != nil {
		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()

		ver, err := h.db.ServerVersion(ctx)
		dbInfo = &DatabaseInfo{
			Engine:    "postgresql",
			Connected: err == nil,
			Version:   ver,
		}
	}

	response.OK(w, r, InfoResponse{
		Service:       h.cfg.AppName,
		Version:       h.cfg.AppVersion,
		Environment:   h.cfg.AppEnv,
		GoVersion:     runtime.Version(),
		OS:            runtime.GOOS,
		Arch:          runtime.GOARCH,
		NumCPU:        runtime.NumCPU(),
		StartedAt:     h.startedAt.Format(time.RFC3339),
		UptimeSeconds: int64(time.Since(h.startedAt).Seconds()),
		Database:      dbInfo,
	})
}
