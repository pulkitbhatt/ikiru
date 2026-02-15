CREATE TABLE monitors (
    id UUID PRIMARY KEY,

    -- ownership
    owner_user_id UUID NOT NULL
        REFERENCES users(id),

    -- identity
    name TEXT NOT NULL,
    description TEXT,

    -- type
    type TEXT NOT NULL
        CHECK (type IN ('http')),

    -- target
    url TEXT NOT NULL,

    -- schedule
    interval_seconds INTEGER NOT NULL
        CHECK (interval_seconds >= 60),

    timeout_ms INTEGER NOT NULL
        CHECK (timeout_ms BETWEEN 100 AND 30000),

    -- lifecycle
    status TEXT NOT NULL
        CHECK (status IN ('active', 'paused')),

    -- soft delete
    deleted_at TIMESTAMPTZ,

    -- checks
    last_checked_at TIMESTAMPTZ,
    next_check_at TIMESTAMPTZ,

    -- audit
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
