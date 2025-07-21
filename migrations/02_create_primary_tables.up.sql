CREATE TABLE IF NOT EXISTS tenants (
                                       id UUID PRIMARY KEY,
                                       name VARCHAR(40) NOT NULL,
                                       created_at TIMESTAMPTZ DEFAULT now(),
                                       owner_id UUID NOT NULL
);

CREATE TABLE IF NOT EXISTS users (
                                     id UUID PRIMARY KEY,
                                     username VARCHAR(20) NOT NULL,
                                     email VARCHAR(254) UNIQUE NOT NULL,
                                     created_at TIMESTAMPTZ DEFAULT now(),
                                     last_seen TIMESTAMPTZ DEFAULT now()
);

CREATE TABLE IF NOT EXISTS events (
    id BIGSERIAL PRIMARY KEY,
    type SMALLINT NOT NULL,
    raid_ref TEXT NULL,
    created_at TIMESTAMPTZ DEFAULT now()
);

CREATE INDEX idx_email_users ON users(email);

CREATE TABLE IF NOT EXISTS storage (
    id UUID PRIMARY KEY
);

CREATE TABLE IF NOT EXISTS user_activities (
    event_id BIGSERIAL REFERENCES events(id) ON DELETE CASCADE,
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    UNIQUE (event_id, user_id)
);

CREATE TABLE IF NOT EXISTS storage_aa_items (
    aa_item_id BIGSERIAL REFERENCES aa_items(id) ON DELETE CASCADE,
    storage_id UUID REFERENCES storage(id) ON DELETE CASCADE,
    since TIMESTAMPTZ DEFAULT now(),
    UNIQUE (aa_item_id, storage_id)
);

CREATE TABLE IF NOT EXISTS publishes (
    id UUID PRIMARY KEY,
    s3 JSONB NOT NULL
)