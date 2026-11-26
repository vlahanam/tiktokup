package server

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"tiktokup/backend/internal/config"
)

// Server wraps http.Server with structured logging and graceful shutdown support.
type Server struct {
	httpServer *http.Server
	cfg        *config.Config
	logger     *slog.Logger
}

// New creates a configured HTTP server instance.
func New(cfg *config.Config, handler http.Handler, logger *slog.Logger) *Server {
	return &Server{
		httpServer: &http.Server{
			Addr:              cfg.Addr(),
			Handler:           handler,
			ReadHeaderTimeout: 10 * time.Second,
			ReadTimeout:       cfg.ReadTimeout,
			WriteTimeout:      cfg.WriteTimeout,
			IdleTimeout:       cfg.IdleTimeout,
		},
		cfg:    cfg,
		logger: logger,
	}
}

// Run starts listening and blocks until ctx is cancelled or the server returns a fatal error.
func (s *Server) Run(ctx context.Context) error {
	errCh := make(chan error, 1)

	go func() {
		s.logger.Info("http_server_starting",
			slog.String("service", s.cfg.AppName),
			slog.String("env", s.cfg.AppEnv),
			slog.String("addr", s.cfg.Addr()),
		)
		if err := s.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
		close(errCh)
	}()

	select {
	case err := <-errCh:
		return err
	case <-ctx.Done():
		s.logger.Info("http_server_shutting_down", slog.Duration("timeout", s.cfg.ShutdownTimeout))
		shutdownCtx, cancel := context.WithTimeout(context.Background(), s.cfg.ShutdownTimeout)
		defer cancel()

		if err := s.httpServer.Shutdown(shutdownCtx); err != nil {
			return err
		}
		s.logger.Info("http_server_stopped")
		return nil
	}
}
