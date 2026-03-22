package model

import (
	"time"

	"github.com/google/uuid"
)

type Incident struct {
	ID           uuid.UUID
	MonitorID    uuid.UUID
	Region       string
	Status       string
	StartedAt    time.Time
	ResolvedAt   time.Time
	FailureCount int
}
