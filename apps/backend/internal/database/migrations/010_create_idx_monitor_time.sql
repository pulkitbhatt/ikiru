CREATE INDEX idx_incidents_monitor_time
ON incidents (monitor_id, started_at DESC);
