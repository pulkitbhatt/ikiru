ALTER TABLE monitor_check_results
ADD CONSTRAINT uniq_monitor_execution
UNIQUE (monitor_id, region, scheduled_at);
