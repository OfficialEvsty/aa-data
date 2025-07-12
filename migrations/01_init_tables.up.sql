CREATE EXTENSION IF NOT EXISTS "pgcrypto";
CREATE TABLE IF NOT EXISTS aa_servers (
    id UUID PRIMARY KEY NOT NULL,
    external_id BIGINT NULL,
    name TEXT NULL
);

CREATE TABLE IF NOT EXISTS aa_nicknames (
    id UUID PRIMARY KEY,
    server_id UUID REFERENCES aa_servers(id) ON DELETE CASCADE,
    name VARCHAR(18) NOT NULL,
    created_at TIMESTAMPTZ DEFAULT now() NOT NULL,
    UNIQUE (server_id, name)
);

CREATE INDEX IF NOT EXISTS idx_nicknames ON aa_nicknames (name);

CREATE TABLE IF NOT EXISTS aa_archived_nicknames (
    id BIGSERIAL PRIMARY KEY,
    nickname_id UUID REFERENCES aa_nicknames(id) ON DELETE CASCADE,
    archived_at TIMESTAMPTZ NOT NULL
);

CREATE TABLE IF NOT EXISTS aa_guilds (
    id UUID PRIMARY KEY NOT NULL,
    name VARCHAR(24) NOT NULL,
    server_id UUID REFERENCES aa_servers(id) ON DELETE CASCADE,
    UNIQUE (name, server_id)
);

CREATE TABLE IF NOT EXISTS aa_bosses (
    id BIGINT PRIMARY KEY,
    name TEXT NOT NULL,
    drop JSONB NOT NULL,
    img_grade_url TEXT NULL,
    img_url TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS aa_items (
    id BIGINT PRIMARY KEY,
    name TEXT NOT NULL,
    tier SMALLINT DEFAULT 0,
    img_grade_url TEXT NULL,
    img_url TEXT NOT NULL
);

