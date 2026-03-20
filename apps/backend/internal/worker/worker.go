package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/pulkitbhatt/ikiru/internal/model"
	"github.com/pulkitbhatt/ikiru/internal/queue"
	"github.com/pulkitbhatt/ikiru/internal/repository"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
)

const (
	MsgReadCount      = 10
	BlockSeconds      = 3
	MinIdleSeconds    = 60
	AutoClaimMsgCount = 50
)

type Worker struct {
	rdb         *redis.Client
	monitorRepo *repository.MonitorRepo
	log         *zerolog.Logger
	stream      string
	workerName  string
	concurrency int
}

func NewWorker(rdb *redis.Client, monitorRepo *repository.MonitorRepo, stream, workerName string, concurrency int, log *zerolog.Logger) *Worker {
	return &Worker{
		rdb:         rdb,
		monitorRepo: monitorRepo,
		log:         log,
		stream:      stream,
		workerName:  workerName,
		concurrency: concurrency,
	}
}

func (w *Worker) Work(ctx context.Context) error {
	w.log.Info().Msg("Starting worker to consume messages...")
	sem := make(chan struct{}, w.concurrency)

	for {
		res, err := w.rdb.XReadGroup(ctx, &redis.XReadGroupArgs{
			Group:    queue.WorkerGroup,
			Consumer: w.workerName,
			Streams:  []string{w.stream, ">"},
			Count:    MsgReadCount,
			Block:    BlockSeconds * time.Second,
		}).Result()

		if err == redis.Nil {
			continue
		}
		if err != nil {
			return fmt.Errorf("error occurred while reading from group %w", err)
		}

		for _, stream := range res {
			for _, msg := range stream.Messages {
				sem <- struct{}{}
				go func(msg redis.XMessage) {
					defer func() { <-sem }()
					w.handleMessage(ctx, msg)
				}(msg)
			}
		}
	}
}

func (w *Worker) handleMessage(ctx context.Context, msg redis.XMessage) {
	start := time.Now()
	job, err := decodeJob(msg)
	if err != nil {
		w.log.Error().Err(err).Msg("failed decoding job")
		return
	}

	if time.Since(job.ScheduledAt) > 2*time.Minute {
		w.log.Info().Str("monitor_id", job.MonitorID.String()).
			Str("job_id", job.JobID).
			Str("scheduled_at", job.ScheduledAt.String()).
			Msg("discarding stale job")
		w.ack(ctx, msg.ID)
		return
	}

	result := ExecuteHTTPCheck(ctx, job)

	res := model.NewMonitorCheckResult(
		job.MonitorID,
		job.Region,
		job.ScheduledAt,
		start,
		time.Now(),
		string(result.Status),
		result.HTTPStatus,
		result.LatencyMs,
		result.Error,
	)

	if err := w.monitorRepo.InsertCheckResult(ctx, res); err != nil {
		w.log.Error().Err(err).Msg("failed to insert monitor check result")
		return
	}

	w.log.Info().Str("monitor_id", job.MonitorID.String()).
		Str("monitor_url", job.URL).
		Str("status", string(result.Status)).
		Int("status_code", result.HTTPStatus).
		Str("error", result.Error).
		Int("latency", result.LatencyMs).
		Msg("checked monitor")

	w.ack(ctx, msg.ID)
}

func (w *Worker) ReclaimPending(ctx context.Context) error {
	start := "0-0"

	for {
		msgs, next, err := w.rdb.XAutoClaim(ctx, &redis.XAutoClaimArgs{
			Stream:   w.stream,
			Group:    queue.WorkerGroup,
			Consumer: w.workerName,
			MinIdle:  MinIdleSeconds * time.Second,
			Start:    start,
			Count:    AutoClaimMsgCount,
		}).Result()

		if err != nil {
			w.log.Error().Err(err).Msg("failed to claim pending messages")
			return fmt.Errorf("error occurred while claiming pending messages: %w", err)
		}

		for _, msg := range msgs {
			w.handleMessage(ctx, msg)
		}

		if next == "0-0" {
			break
		}

		start = next
	}

	return nil
}

func (w *Worker) ack(ctx context.Context, id string) {
	err := w.rdb.XAck(ctx, w.stream, queue.WorkerGroup, id).Err()
	if err != nil {
		w.log.Error().Err(err).Msg("ack failed")
	}
}

func decodeJob(msg redis.XMessage) (queue.MonitorJob, error) {
	var job queue.MonitorJob
	payload := msg.Values["payload"].(string)
	err := json.Unmarshal([]byte(payload), &job)
	return job, err
}
