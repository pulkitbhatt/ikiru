CREATE UNIQUE INDEX uniq_open_incident
ON incidents (monitor_id, region)
WHERE status = 'open';
