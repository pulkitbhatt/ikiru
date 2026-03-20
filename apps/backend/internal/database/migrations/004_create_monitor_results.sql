CREATE TABLE monitor_check_results (
    id UUID PRIMARY KEY,

    monitor_id UUID NOT NULL
        REFERENCES monitors(id) ON DELETE CASCADE,

    region TEXT NOT NULL,

    scheduled_at TIMESTAMPTZ NOT NULL,

    started_at TIMESTAMPTZ NOT NULL,
    finished_at TIMESTAMPTZ NOT NULL,

    status TEXT NOT NULL
        CHECK (status IN ('success', 'failure', 'timeout')),

    http_status INTEGER,
    latency_ms INTEGER,

    error TEXT,

    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
