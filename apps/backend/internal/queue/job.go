package queue

import (
	"time"

	"github.com/google/uuid"
)

type MonitorJob struct {
	JobID       string
	MonitorID   uuid.UUID
	Region      string
	URL         string
	TimeoutMs   int
	ScheduledAt time.Time
}
