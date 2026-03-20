package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/pulkitbhatt/ikiru/internal/config"
	"github.com/pulkitbhatt/ikiru/internal/database"
	"github.com/pulkitbhatt/ikiru/internal/logger"
	"github.com/pulkitbhatt/ikiru/internal/queue"
	"github.com/pulkitbhatt/ikiru/internal/repository"
	"github.com/pulkitbhatt/ikiru/internal/worker"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg, err := config.LoadConfig()
	if err != nil {
		panic("failed to load config" + err.Error())
	}

	logger := logger.New(cfg)
	db, err := database.New(cfg, &logger)
	if err != nil {
		logger.Fatal().Msg("failed to initialize database")
	}
	defer db.Pool.Close()

	repo := repository.NewMonitorRepo(db.Pool)
	redisClient := queue.NewRedis(cfg.Redis.Address)
	workerName := fmt.Sprintf("worker-%d", os.Getpid())

	if err := queue.EnsureConsumerGroup(ctx, redisClient, cfg.Redis.Stream); err != nil {
		panic(err)
	}

	worker := worker.NewWorker(
		redisClient,
		repo,
		cfg.Redis.Stream,
		workerName,
		config.WorkerMaxConcurrency,
		&logger,
	)

	err = worker.ReclaimPending(ctx)
	if err != nil {
		panic(err)
	}

	if err := worker.Work(ctx); err != nil {
		panic(err)
	}
}
