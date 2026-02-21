package dto

import "github.com/google/uuid"

type CreateMonitorRequest struct {
	Name            string  `json:"name"`
	Description     *string `json:"description"`
	URL             string  `json:"url"`
	IntervalSeconds int     `json:"interval_seconds"`
	TimeoutMs       int     `json:"timeout_ms"`
}

type DueMonitor struct {
	ID        uuid.UUID `json:"id"`
	UserId    string    `json:"user_id"`
	Type      string    `json:"type"`
	URL       string    `json:"url"`
	TimeoutMs int       `json:"timeout_ms"`
}
