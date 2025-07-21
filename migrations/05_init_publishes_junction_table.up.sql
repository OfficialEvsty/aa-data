

CREATE TABLE IF NOT EXISTS tenant_publishes (
                                                  tenant_id UUID REFERENCES tenants(id) ON DELETE CASCADE,
                                                  publish_id UUID REFERENCES publishes(id) ON DELETE CASCADE,
                                                  user_id UUID REFERENCES users(id) ON DELETE CASCADE,
                                                  published_at TIMESTAMPTZ DEFAULT now(),
                                                  UNIQUE (tenant_id, publish_id)
);

CREATE TABLE IF NOT EXISTS aa_tenant_guilds (
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    guild_id UUID NOT NULL REFERENCES aa_guilds(id) ON DELETE CASCADE,
    joined_at TIMESTAMPTZ DEFAULT now()
);
