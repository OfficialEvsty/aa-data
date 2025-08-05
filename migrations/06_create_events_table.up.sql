CREATE TABLE IF NOT EXISTS aa_template_events (
    id BIGSERIAL PRIMARY KEY,
    name TEXT NULL,
    metadata JSONB NULL
);

CREATE TABLE IF NOT EXISTS aa_events (
    id BIGSERIAL PRIMARY KEY,
    template_id BIGINT REFERENCES aa_template_events(id) ON DELETE CASCADE,
    occurred_at TIMESTAMPTZ DEFAULT now()
);

DROP TABLE IF EXISTS aa_guild_storage;
DROP TABLE IF EXISTS aa_guild_events;
DROP TABLE IF EXISTS aa_guild_users;


CREATE TABLE IF NOT EXISTS aa_tenant_storage (
                                                tenant_id UUID REFERENCES tenants(id) ON DELETE CASCADE,
                                                storage_id UUID REFERENCES storage(id) ON DELETE CASCADE,
                                                UNIQUE(tenant_id, storage_id)
);

CREATE TABLE IF NOT EXISTS aa_event_bosses (
    event_template_id INTEGER REFERENCES aa_template_events(id) ON DELETE CASCADE,
    boss_id BIGINT REFERENCES aa_bosses(id) ON DELETE CASCADE,
    UNIQUE (event_template_id, boss_id)
);

