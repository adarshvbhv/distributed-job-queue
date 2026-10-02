package http

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/adarshvbhv/distributed-job-queue/internal/config"
	"github.com/adarshvbhv/distributed-job-queue/internal/storage"
)

type Server struct {
	httpServer *http.Server
	logger     *slog.Logger
	config     *config.Config
	db         *storage.Postgres
}

func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "OK")
}

func NewServer(logger *slog.Logger, config *config.Config, db *storage.Postgres) *Server {

	mux := http.NewServeMux()

	s := &Server{
		httpServer: &http.Server{
			Addr:    ":" + config.HTTPPort,
			Handler: mux,
		},
		logger: logger,
		config: config,
		db:     db,
	}

	s.requestHandler(mux)
	return s

}

func (s *Server) Start() error {

	s.logger.Info("starting server on", "port", s.config.HTTPPort)
	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {

	return s.httpServer.Shutdown(ctx)

}

func (s *Server) requestHandler(mux *http.ServeMux) {

	mux.HandleFunc("GET /health", s.healthHandler)

}
