CREATE INDEX idx_results_monitor_region_time
ON monitor_check_results (monitor_id, region, scheduled_at DESC);
