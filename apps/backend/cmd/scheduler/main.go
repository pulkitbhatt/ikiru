package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/pulkitbhatt/ikiru/internal/config"
	"github.com/pulkitbhatt/ikiru/internal/database"
	"github.com/pulkitbhatt/ikiru/internal/logger"
	"github.com/pulkitbhatt/ikiru/internal/publisher"
	"github.com/pulkitbhatt/ikiru/internal/repository"
	"github.com/pulkitbhatt/ikiru/internal/scheduler"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg, err := config.LoadConfig()
	if err != nil {
		panic("failed to load config" + err.Error())
	}

	logger := logger.New(cfg)
	logger.Info().Msg("Starting Scheduler...")

	db, err := database.New(cfg, &logger)
	if err != nil {
		logger.Fatal().Msg("failed to initialize database")
	}
	defer db.Pool.Close()

	repo := repository.NewMonitorRepo(db.Pool)
	pub := publisher.FakePublisher{}
	scheduler := scheduler.New(&pub, repo, &logger)
	scheduler.Run(ctx)
}
