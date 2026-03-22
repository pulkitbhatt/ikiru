package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/pulkitbhatt/ikiru/internal/config"
	"github.com/pulkitbhatt/ikiru/internal/database"
	eventrouter "github.com/pulkitbhatt/ikiru/internal/event_router"
	"github.com/pulkitbhatt/ikiru/internal/logger"
	"github.com/pulkitbhatt/ikiru/internal/outbox"
	"github.com/pulkitbhatt/ikiru/internal/publisher"
	"github.com/pulkitbhatt/ikiru/internal/queue"
	"github.com/pulkitbhatt/ikiru/internal/repository"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}

	logger := logger.New(cfg)
	db, err := database.New(cfg, &logger)
	if err != nil {
		panic(err)
	}

	outboxRepo := repository.NewOutboxRepo(db.Pool)
	rdb := queue.NewRedis(cfg.Redis.Address)
	redisPublisher := publisher.NewRedisPublisher(rdb, &logger)
	router := eventrouter.New(redisPublisher)
	outboxPublisher := outbox.NewOutboxPublisher(
		outboxRepo,
		router,
		&logger,
	)

	outboxPublisher.Run(ctx)
}
