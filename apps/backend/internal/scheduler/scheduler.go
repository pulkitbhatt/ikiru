package scheduler

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/pulkitbhatt/ikiru/internal/config"
	"github.com/pulkitbhatt/ikiru/internal/model"
	"github.com/pulkitbhatt/ikiru/internal/publisher"
	"github.com/pulkitbhatt/ikiru/internal/queue"
	"github.com/pulkitbhatt/ikiru/internal/util"
	"github.com/rs/zerolog"
)

const (
	MonitorBatchLimit = 100
	PollInterval      = 10
)

type MonitorRepo interface {
	ClaimDueMonitors(context.Context, int) ([]model.Monitor, error)
}

type Scheduler struct {
	publisher   publisher.Publisher
	monitorRepo MonitorRepo
	logger      *zerolog.Logger
}

func New(pub publisher.Publisher, monitorRepo MonitorRepo, logger *zerolog.Logger) *Scheduler {
	return &Scheduler{
		publisher:   pub,
		monitorRepo: monitorRepo,
		logger:      logger,
	}
}

func (s *Scheduler) Run(ctx context.Context) {
	ticker := time.NewTicker(PollInterval * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			s.logger.Info().Msg("scheduler shutting down")
			return
		case <-ticker.C:
			claimedMonitors, err := s.monitorRepo.ClaimDueMonitors(ctx, MonitorBatchLimit)
			if err != nil {
				s.logger.Error().Err(err).Msg("failed to claim monitors")
				continue
			}

			s.logger.Info().Msgf("Claimed %v monitors", len(claimedMonitors))

			for _, m := range claimedMonitors {
				s.logger.Info().
					Str("monitor_id", m.ID.String()).
					Msg("scheduling job for monitor")
				for _, region := range config.RegionsToMonitor {
					job := queue.MonitorJob{
						JobID:       util.GenerateUUIDStr(),
						MonitorID:   m.ID,
						Region:      region,
						URL:         m.URL,
						TimeoutMs:   m.TimeoutMs,
						ScheduledAt: time.Now(),
					}
					stream := fmt.Sprintf("monitor_checks:%s", region)
					payload, _ := json.Marshal(job)
					msg := publisher.Message{
						Destination: stream,
						Payload:     payload,
					}
					if err := s.publisher.Publish(ctx, msg); err != nil {
						s.logger.Error().
							Err(err).
							Str("job_id", job.JobID).
							Str("monitor_id", job.MonitorID.String()).
							Str("region", job.Region).
							Msg("failed to publish job")
						continue
					}
					s.logger.Info().Str("job_id", job.JobID).
						Str("scheduled_at", job.ScheduledAt.String()).
						Msg("Published job successfully")
				}
			}
		}
	}
}
