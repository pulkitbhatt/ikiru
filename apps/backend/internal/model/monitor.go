package model

import (
	"fmt"
	"math/rand"
	"net/url"
	"slices"
	"time"

	"github.com/google/uuid"

	"github.com/pulkitbhatt/ikiru/internal/util"
	"github.com/pulkitbhatt/ikiru/internal/validation"
)

type MonitorType string

const (
	MonitorTypeHTTP MonitorType = "http"
)

type MonitorStatus string

const (
	MonitorStatusActive MonitorStatus = "active"
	MonitorStatusPaused MonitorStatus = "paused"
)

var AllowedIntervals = []int{60, 120, 300, 600, 1800, 3600}
var AllowedMonitorTypes = []MonitorType{MonitorTypeHTTP}
var AllowedStatuses = []MonitorStatus{
	MonitorStatusActive, MonitorStatusPaused,
}

const Jitter = 20

type Monitor struct {
	ID uuid.UUID

	OwnerUserID uuid.UUID

	Name        string
	Description *string

	Type MonitorType

	URL string

	IntervalSeconds int
	TimeoutMs       int

	Status    MonitorStatus
	DeletedAt *time.Time

	LastCheckedAt *time.Time
	NextCheckAt   time.Time

	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewMonitor(
	ownerUserID uuid.UUID,
	name string,
	url string,
	intervalSeconds int,
	timeoutMs int,
	description *string,
) *Monitor {

	now := time.Now().UTC()

	return &Monitor{
		ID:              util.GenerateUUID(),
		OwnerUserID:     ownerUserID,
		Name:            name,
		Description:     description,
		Type:            MonitorTypeHTTP,
		URL:             url,
		IntervalSeconds: intervalSeconds,
		TimeoutMs:       timeoutMs,

		Status:    MonitorStatusActive,
		DeletedAt: nil,

		LastCheckedAt: nil,
		NextCheckAt: now.
			Add(time.Duration(intervalSeconds) * time.Second).
			Add(randomJitter(Jitter)),
	}
}

func (m *Monitor) Validate() error {
	var errs validation.ValidationErrors

	if m.Name == "" {
		errs = errs.Add("name", "name is required")
	}

	if m.URL == "" {
		errs = errs.Add("url", "url is required")
	} else if u, err := url.ParseRequestURI(m.URL); err != nil || (u.Scheme != "http" && u.Scheme != "https") {
		errs = errs.Add("url", "invalid URL, must be http or https")
	}

	if !slices.Contains(AllowedIntervals, m.IntervalSeconds) {
		errs = errs.Add("interval_seconds", fmt.Sprintf(
			"must be one of the allowed values: %v", AllowedIntervals))
	}

	if m.TimeoutMs < 100 || m.TimeoutMs > 30000 {
		errs = errs.Add("timeout_ms", "timeout ms must be between 100 and 30000")
	}

	if !m.Type.IsValid() {
		errs = errs.Add("type", fmt.Sprintf(
			"invalid monitor type: %q, allowed values: %v", m.Type, AllowedMonitorTypes))
	}

	if m.OwnerUserID == uuid.Nil {
		errs = errs.Add("owner_user_id", "owner user id is required")
	}

	if !m.Status.IsValid() {
		errs = errs.Add("status", fmt.Sprintf(
			"invalid monitor status: %q, allowed values: %v", m.Status, AllowedStatuses))
	}

	if errs.HasErrors() {
		return errs
	}

	return nil
}

func (t MonitorType) IsValid() bool {
	return slices.Contains(AllowedMonitorTypes, t)
}

func (s MonitorStatus) IsValid() bool {
	return slices.Contains(AllowedStatuses, s)
}

func randomJitter(n int) time.Duration {
	jitter := rand.Intn(n)
	return time.Duration(jitter) * time.Second
}
