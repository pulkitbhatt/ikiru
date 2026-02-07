BEGIN;

CREATE TABLE users (
    id UUID PRIMARY KEY,

    -- External identity provider user identifier
    idp_user_id TEXT NOT NULL UNIQUE,

    email TEXT,

    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

COMMIT;
