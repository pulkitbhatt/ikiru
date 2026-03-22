package publisher

import (
	"context"

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

func (p *RedisPublisher) Publish(ctx context.Context, msg Message) error {
	return p.rdb.XAdd(ctx, &redis.XAddArgs{
		Stream: msg.Destination,
		Values: map[string]any{
			"payload": msg.Payload,
		},
		MaxLen: StreamMaxLen,
		Approx: true,
	}).Err()
}
