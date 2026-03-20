package publisher

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/pulkitbhatt/ikiru/internal/queue"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
)

const (
	StreamMaxLen = 100000
)

type RedisPublisher struct {
	rdb    *redis.Client
	logger *zerolog.Logger
}

func NewRedisPublisher(rdb *redis.Client, logger *zerolog.Logger) *RedisPublisher {
	return &RedisPublisher{
		rdb:    rdb,
		logger: logger,
	}
}

func (p *RedisPublisher) Publish(ctx context.Context, job queue.MonitorJob) error {
	stream := fmt.Sprintf("monitor_checks:%s", job.Region)

	payload, err := json.Marshal(job)
	if err != nil {
		return err
	}

	return p.rdb.XAdd(ctx, &redis.XAddArgs{
		Stream: stream,
		Values: map[string]any{
			"payload": payload,
		},
		MaxLen: StreamMaxLen,
		Approx: true,
	}).Err()
}
