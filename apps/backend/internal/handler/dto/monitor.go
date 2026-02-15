package dto

type CreateMonitorRequest struct {
	Name            string  `json:"name"`
	Description     *string `json:"description"`
	URL             string  `json:"url"`
	IntervalSeconds int     `json:"interval_seconds"`
	TimeoutMs       int     `json:"timeout_ms"`
}
