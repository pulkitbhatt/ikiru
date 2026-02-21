package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/pulkitbhatt/ikiru/internal/config"
	"github.com/pulkitbhatt/ikiru/internal/database"
	"github.com/pulkitbhatt/ikiru/internal/handler"
	"github.com/pulkitbhatt/ikiru/internal/logger"
	"github.com/pulkitbhatt/ikiru/internal/repository"
	"github.com/pulkitbhatt/ikiru/internal/router"
	"github.com/pulkitbhatt/ikiru/internal/server"
	"github.com/pulkitbhatt/ikiru/internal/service"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		panic("failed to load config" + err.Error())
	}

	logger := logger.New(cfg)
	logger.Info().Msg("Starting server...")

	if err := database.Migrate(context.Background(), &logger, cfg); err != nil {
		logger.Fatal().Err(err).Msg("failed to migrate database")
	}

	svr, err := server.New(cfg, &logger)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to initialize server")
	}

	repos := repository.NewRepositories(svr)
	services := service.NewServices(svr, repos)
	handlers := handler.NewHandlers(svr, services)

	router := router.NewRouter(svr, handlers)
	svr.Setup(router)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	go func() {
		if err := svr.Start(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Fatal().Err(err).Msg("failed to start server")
		}
	}()

	<-ctx.Done()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := svr.Stop(ctx); err != nil {
		logger.Fatal().Err(err).Msg("failed to stop server")
	}

	logger.Info().Msg("Server stopped gracefulluy")
}
