package server

import (
	"booking/internal/config"
	"context"
	"errors"
	"net/http"

	"golang.org/x/exp/slog"
)

type Server struct {
	HttpServer *http.Server
	log        *slog.Logger
}

func New(cfg *config.Config, router http.Handler, log *slog.Logger) *Server {
	return &Server{
		HttpServer: &http.Server{
			Addr:         cfg.Address,
			Handler:      router,
			ReadTimeout:  cfg.Timeout,
			WriteTimeout: cfg.Timeout,
			IdleTimeout:  cfg.IdleTimeout,
		},
		log: log,
	}
}

func (s *Server) Run() error {
	s.log.Info("server started", "addr", s.HttpServer.Addr)

	err := s.HttpServer.ListenAndServe()

	if errors.Is(err, http.ErrServerClosed) {
		return nil
	}

	return err
}

func (s *Server) Shutdown(ctx context.Context) error {
	s.log.Info("stopping server")

	return s.HttpServer.Shutdown(ctx)
}
