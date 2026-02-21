package scheduler

import (
	"context"
	"time"

	"github.com/pulkitbhatt/ikiru/internal/handler/dto"
	"github.com/rs/zerolog"
)

const (
	MonitorBatchLimit = 100
	PollInterval      = 10
)

type Publisher interface {
	Publish(ctx context.Context, payload dto.DueMonitor)
}

type MonitorRepo interface {
	ClaimDueMonitors(context.Context, int) ([]dto.DueMonitor, error)
}

type Scheduler struct {
	publisher   Publisher
	monitorRepo MonitorRepo
	logger      *zerolog.Logger
}

func New(pub Publisher, monitorRepo MonitorRepo, logger *zerolog.Logger) *Scheduler {
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

			for _, m := range claimedMonitors {
				s.logger.Info().Str("monitor_id", m.ID.String()).Msg("scheduling job for monitor")
				// for testing, will just log to console...
				go s.publisher.Publish(ctx, m)
			}
		}
	}
}
