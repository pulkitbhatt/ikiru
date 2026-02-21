CREATE INDEX idx_monitors_scheduler
ON monitors (next_check_at, id)
WHERE status = 'active'
  AND deleted_at IS NULL;
