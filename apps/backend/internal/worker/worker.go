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

	// jobs older than this would be considered stale
	StaleJobThreshold = 2
	FailureThreshold  = 3
	ResolveThreshold  = 2
)

type Worker struct {
	rdb              *redis.Client
	monitorCheckRepo *repository.MonitorCheckRepo
	incidentRepo     *repository.IncidentRepo

	log *zerolog.Logger

	stream      string
	workerName  string
	concurrency int
}

func NewWorker(rdb *redis.Client,
	monitorCheckRepo *repository.MonitorCheckRepo,
	incidentRepo *repository.IncidentRepo,
	stream, workerName string,
	concurrency int,
	log *zerolog.Logger,
) *Worker {
	return &Worker{
		rdb:              rdb,
		monitorCheckRepo: monitorCheckRepo,
		incidentRepo:     incidentRepo,
		log:              log,
		stream:           stream,
		workerName:       workerName,
		concurrency:      concurrency,
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
		w.log.Error().Err(err).
			Str("monitor_id", job.MonitorID.String()).
			Str("region", job.Region).
			Msg("failed decoding job")
		return
	}

	if isStaleJob(job.ScheduledAt) {
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

	if err := w.monitorCheckRepo.InsertCheckResult(ctx, res); err != nil {
		w.log.Error().Err(err).
			Str("monitor_id", job.MonitorID.String()).
			Str("region", job.Region).
			Msg("failed to insert monitor check result")
		return
	}

	w.log.Info().Str("monitor_id", job.MonitorID.String()).
		Str("monitor_url", job.URL).
		Str("status", string(result.Status)).
		Int("status_code", result.HTTPStatus).
		Str("error", result.Error).
		Int("latency", result.LatencyMs).
		Msg("checked monitor")

	w.evaluateIncident(ctx, job)

	w.ack(ctx, msg.ID)
}

func (w *Worker) evaluateIncident(ctx context.Context, job queue.MonitorJob) {
	results, err := w.monitorCheckRepo.GetLastNResults(ctx,
		job.MonitorID,
		job.Region,
		max(FailureThreshold, ResolveThreshold),
	)
	if err != nil {
		w.log.Error().Err(err).
			Str("monitor_id", job.MonitorID.String()).
			Str("region", job.Region).
			Msg("failed to fetch last n checks")
		return
	}

	allFailures := hasNFailures(results, FailureThreshold)

	openIncident, err := w.incidentRepo.GetOpenIncident(ctx,
		job.MonitorID,
		job.Region,
	)
	if err != nil {
		w.log.Error().Err(err).
			Str("monitor_id", job.MonitorID.String()).
			Str("region", job.Region).
			Msg("failed to fetch open incidents")
		return
	}

	if allFailures {
		if openIncident == nil {
			ok, err := w.incidentRepo.TryCreateIncidentWithOutbox(ctx,
				job.MonitorID,
				job.Region,
				FailureThreshold,
			)
			if err != nil {
				w.log.Error().Err(err).
					Str("monitor_id", job.MonitorID.String()).
					Str("region", job.Region).
					Msg("failed to create new incident")
				return
			}

			if ok {
				w.log.Info().
					Str("monitor_id", job.MonitorID.String()).
					Str("region", job.Region).
					Msg("Created new incident successfully")
			}
		} else if isFailure(results[0]) {
			if err := w.incidentRepo.IncrementFailureCount(ctx,
				job.MonitorID,
				job.Region,
			); err != nil {
				w.log.Error().Err(err).
					Str("monitor_id", job.MonitorID.String()).
					Str("region", job.Region).
					Msg("failed to increment failure count")
				return
			}
			w.log.Debug().
				Str("monitor_id", job.MonitorID.String()).
				Str("region", job.Region).
				Msg("Incremented failure count successfully")
		}
		return
	}

	if openIncident == nil {
		return
	}

	allSuccess := hasNSuccess(results, ResolveThreshold)
	if !allSuccess {
		return
	}

	ok, err := w.incidentRepo.ResolveIncidentWithOutbox(ctx,
		job.MonitorID,
		job.Region,
	)
	if err != nil {
		w.log.Error().Err(err).Msg("failed to resolve incident")
		return
	}
	if ok {
		w.log.Info().
			Str("monitor_id", job.MonitorID.String()).
			Str("region", job.Region).
			Msg("Resolved incident successfully")
	}
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

func isStaleJob(scheduledAt time.Time) bool {
	return time.Since(scheduledAt) > (StaleJobThreshold * time.Minute)
}

func hasNFailures(results []string, n int) bool {
	if len(results) < n {
		return false
	}
	for i := range n {
		if !isFailure(results[i]) {
			return false
		}
	}
	return true
}

func hasNSuccess(results []string, n int) bool {
	if len(results) < n {
		return false
	}
	for i := range n {
		if results[i] != "success" {
			return false
		}
	}
	return true
}

func isFailure(res string) bool {
	return res == "failure" || res == "timeout"
}
