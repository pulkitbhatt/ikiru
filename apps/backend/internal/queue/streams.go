package queue

import (
	"context"
	"strings"

	"github.com/redis/go-redis/v9"
)

const WorkerGroup = "http-workers"

func EnsureConsumerGroup(
	ctx context.Context,
	rdb *redis.Client,
	stream string,
) error {

	err := rdb.XGroupCreateMkStream(
		ctx,
		stream,
		WorkerGroup,
		"$",
	).Err()

	if err != nil && !isBusyGroupErr(err) {
		return err
	}

	return nil
}

func isBusyGroupErr(err error) bool {
	if err == nil {
		return false
	}
	return strings.HasPrefix(err.Error(), "BUSYGROUP")
}
