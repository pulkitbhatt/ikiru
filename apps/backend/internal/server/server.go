package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/pulkitbhatt/ikiru/internal/config"
	"github.com/pulkitbhatt/ikiru/internal/database"
	"github.com/rs/zerolog"
)

type Server struct {
	Config     *config.Config
	Logger     *zerolog.Logger
	httpServer *http.Server
	Db         *database.Database
}

func New(cfg *config.Config, logger *zerolog.Logger) (*Server, error) {
	db, err := database.New(cfg, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}

	return &Server{
		Config: cfg,
		Logger: logger,
		Db:     db,
	}, nil
}

func (s *Server) Setup(h http.Handler) {
	s.httpServer = &http.Server{
		Addr:    ":" + s.Config.Server.Port,
		Handler: h,
	}
}

func (s *Server) Start() error {
	if s.httpServer == nil {
		return errors.New("HTTP server not initialized")
	}

	s.Logger.Info().Msg("Starting server on port " + s.Config.Server.Port)
	return s.httpServer.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	if err := s.httpServer.Shutdown(ctx); err != nil {
		return fmt.Errorf("failed to stop HTTP server: %w", err)
	}

	s.Logger.Info().Msg("Server stopped")
	return nil
}
