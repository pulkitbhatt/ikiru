package eventrouter

import (
	"context"
	"fmt"

	"github.com/pulkitbhatt/ikiru/internal/config"
	"github.com/pulkitbhatt/ikiru/internal/model"
	"github.com/pulkitbhatt/ikiru/internal/publisher"
)

type Route struct {
	Publishers  []publisher.Publisher
	Destination string
}

type EventRouter struct {
	routes map[string]Route
}

func New(redisPublisher *publisher.RedisPublisher) *EventRouter {
	return &EventRouter{
		routes: map[string]Route{
			config.EventIncidentCreated: {
				Publishers:  []publisher.Publisher{redisPublisher},
				Destination: "incident_events",
			},
			config.EventIncidentResolved: {
				Publishers:  []publisher.Publisher{redisPublisher},
				Destination: "incident_events",
			},
		},
	}
}

func (r *EventRouter) Route(ctx context.Context, e model.OutboxEvent) error {
	route, ok := r.routes[e.Type]
	if !ok {
		return nil
	}

	for _, pub := range route.Publishers {
		msg := publisher.Message{
			Destination: route.Destination,
			Payload:     e.Payload,
		}
		if err := pub.Publish(ctx, msg); err != nil {
			return fmt.Errorf(
				"publish failed: event_type=%s destination=%s: %w",
				e.Type,
				msg.Destination,
				err,
			)
		}
	}
	return nil
}
