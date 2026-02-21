package scheduler

import "github.com/google/uuid"

type DueMonitor struct {
	ID              uuid.UUID
	URL             string
	IntervalSeconds int
	TimeoutMs       int
}
