package outbox

import (
	"context"
	"time"

	eventrouter "github.com/pulkitbhatt/ikiru/internal/event_router"
	"github.com/pulkitbhatt/ikiru/internal/repository"
	"github.com/rs/zerolog"
)

const (
	PollIntervalSeconds = 10
	BatchSize           = 100
)

type OutboxPublisher struct {
	outboxRepo *repository.OutboxRepo
	router     *eventrouter.EventRouter
	log        *zerolog.Logger
}

func NewOutboxPublisher(
	outboxRepo *repository.OutboxRepo,
	router *eventrouter.EventRouter,
	log *zerolog.Logger) *OutboxPublisher {
	return &OutboxPublisher{
		outboxRepo: outboxRepo,
		router:     router,
		log:        log,
	}
}

func (p *OutboxPublisher) Run(ctx context.Context) {
	p.log.Info().Msg("Outbox publisher started...")
	ticker := time.NewTicker(PollIntervalSeconds * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			p.log.Info().Msg("outbox publisher shutting down")
			return

		case <-ticker.C:
			p.processBatch(ctx)
		}
	}
}

func (p *OutboxPublisher) processBatch(ctx context.Context) {
	events, err := p.outboxRepo.FetchUnprocessed(ctx, BatchSize)
	if err != nil {
		p.log.Error().Err(err).Msg("failed to fetch unprocessed events")
		return
	}

	for _, e := range events {
		if err := p.router.Route(ctx, e); err != nil {
			p.log.Error().Err(err).
				Str("event_id", e.ID.String()).
				Str("event_type", e.Type).
				Msg("failed to route event")
			return
		}

		p.log.Info().
			Str("event_id", e.ID.String()).
			Str("event_type", e.Type).
			Msg("Routed event successfully")

		if err := p.outboxRepo.MarkProcessed(ctx, e.ID); err != nil {
			p.log.Error().Err(err).
				Str("event_id", e.ID.String()).
				Str("event_type", e.Type).
				Msg("failed to mark event as processed")
			return
		}

		p.log.Info().
			Str("event_id", e.ID.String()).
			Str("event_type", e.Type).
			Msg("Marked event as processed")
	}
}
