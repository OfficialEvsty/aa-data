

CREATE TABLE IF NOT EXISTS tokens (
    token TEXT NOT NULL,
    user_id UUID PRIMARY KEY,
    expires_at TIMESTAMPTZ NOT NULL
)