CREATE INDEX idx_outbox_ready
ON outbox_events (created_at)
WHERE processed_at IS NULL AND processing_at IS NULL;
