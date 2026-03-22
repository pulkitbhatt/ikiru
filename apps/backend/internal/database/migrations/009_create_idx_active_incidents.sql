CREATE INDEX idx_incidents_active
ON incidents (monitor_id, region)
WHERE status = 'open';
