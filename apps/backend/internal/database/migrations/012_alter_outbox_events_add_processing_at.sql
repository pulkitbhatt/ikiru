ALTER TABLE outbox_events
ADD COLUMN processing_at TIMESTAMPTZ;
