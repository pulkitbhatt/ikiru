CREATE TABLE incidents (
    id UUID PRIMARY KEY,

    monitor_id UUID NOT NULL
        REFERENCES monitors(id) ON DELETE CASCADE,

    region TEXT NOT NULL,

    status TEXT NOT NULL
        CHECK (status IN ('open', 'resolved')),

    started_at TIMESTAMPTZ NOT NULL,
    resolved_at TIMESTAMPTZ,

    failure_count INTEGER DEFAULT 0,

    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
