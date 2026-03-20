package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/pulkitbhatt/ikiru/internal/util"
)

type MonitorCheckResult struct {
	ID          uuid.UUID
	MonitorID   uuid.UUID
	Region      string
	ScheduledAt time.Time
	StartedAt   time.Time
	FinishedAt  time.Time
	Status      string
	HTTPStatus  int
	LatencyMs   int
	Error       string
}

func NewMonitorCheckResult(
	monitorID uuid.UUID,
	region string,
	scheduledAt time.Time,
	startedAt time.Time,
	finishedAt time.Time,
	status string,
	httpStatus int,
	latencyMs int,
	error string,
) MonitorCheckResult {
	return MonitorCheckResult{
		ID:          util.GenerateUUID(),
		MonitorID:   monitorID,
		Region:      region,
		ScheduledAt: scheduledAt,
		FinishedAt:  finishedAt,
		Status:      status,
		HTTPStatus:  httpStatus,
		LatencyMs:   latencyMs,
		Error:       error,
	}
}
