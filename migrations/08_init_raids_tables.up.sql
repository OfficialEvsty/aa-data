CREATE TABLE IF NOT EXISTS raids (
    id UUID PRIMARY KEY,
    publish_id UUID UNIQUE REFERENCES publishes(id) ON DELETE CASCADE,
    raid_at TIMESTAMPTZ NULL,
    created_at TIMESTAMPTZ DEFAULT now(),
    attendance INTEGER DEFAULT 0,
    status TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS raid_events (
    raid_id UUID REFERENCES raids(id) ON DELETE CASCADE,
    event_id INTEGER REFERENCES aa_template_events(id) ON DELETE CASCADE,
    UNIQUE(raid_id, event_id)
);

CREATE TABLE IF NOT EXISTS lunark (
    id UUID PRIMARY KEY,
    name TEXT NULL,
    start_date TIMESTAMPTZ NOT NULL,
    end_date TIMESTAMPTZ NULL,
    opened BOOLEAN DEFAULT TRUE
);

CREATE TABLE IF NOT EXISTS lunark_raids (
    lunark_id UUID REFERENCES lunark(id) ON DELETE CASCADE,
    raid_id UUID REFERENCES raids(id) ON DELETE CASCADE,
    UNIQUE(lunark_id, raid_id)
);

CREATE TABLE IF NOT EXISTS attendance (
    raid_id UUID REFERENCES raids(id) ON DELETE CASCADE,
    nickname_id UUID REFERENCES aa_nicknames(id) ON DELETE CASCADE,
    UNIQUE(raid_id, nickname_id)
);

CREATE TABLE IF NOT EXISTS raid_items (
    raid_id UUID REFERENCES raids(id) ON DELETE CASCADE,
    item_id BIGINT REFERENCES aa_items(id) ON DELETE CASCADE,
    rate BIGINT NOT NULL,
    UNIQUE(raid_id, item_id)
);

CREATE TABLE IF NOT EXISTS tenant_lunark (
    tenant_id UUID REFERENCES tenants(id) ON DELETE CASCADE,
    lunark_id UUID REFERENCES lunark(id) ON DELETE CASCADE,
    UNIQUE(lunark_id)
)
