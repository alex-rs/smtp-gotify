package health

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
)

type Response struct {
	Status string `json:"status"`
}

type Server struct {
	server *http.Server
	logger *slog.Logger
}

func NewServer(addr string, logger *slog.Logger) *Server {
	mux := http.NewServeMux()
	s := &Server{
		server: &http.Server{
			Addr:    addr,
			Handler: mux,
		},
		logger: logger,
	}

	mux.HandleFunc("/health", s.handleHealth)

	return s
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(Response{Status: "ok"}); err != nil {
		s.logger.Error("failed to encode health response", "error", err)
	}
}

func (s *Server) ListenAndServe() error {
	s.logger.Info("starting health server", "addr", s.server.Addr)
	return s.server.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	s.logger.Info("stopping health server")
	return s.server.Shutdown(ctx)
}
