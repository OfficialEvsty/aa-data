CREATE TABLE IF NOT EXISTS user_aa_nicknames (
                                                 user_id UUID REFERENCES users(id) ON DELETE CASCADE,
                                                 nickname_id UUID REFERENCES aa_nicknames(id) ON DELETE CASCADE,
                                                 is_archived BOOLEAN DEFAULT FALSE,
                                                 UNIQUE (user_id, nickname_id)
);

CREATE TABLE IF NOT EXISTS event_type_bosses (
                                                 event_type SMALLINT NOT NULL,
                                                 boss_id BIGINT REFERENCES aa_bosses(id) ON DELETE CASCADE,
                                                 UNIQUE (event_type, boss_id)
);

CREATE TABLE IF NOT EXISTS aa_guild_users (
                                              guild_id UUID REFERENCES aa_guilds(id) ON DELETE CASCADE,
                                              user_id UUID REFERENCES users(id) ON DELETE CASCADE,
                                              since TIMESTAMPTZ DEFAULT now(),
                                              UNIQUE (guild_id, user_id)
);


CREATE TABLE IF NOT EXISTS aa_guild_nicknames (
                                                  guild_id UUID REFERENCES aa_guilds(id) ON DELETE CASCADE,
                                                  nickname_id UUID REFERENCES aa_nicknames(id) ON DELETE CASCADE,
                                                  member_at TIMESTAMPTZ DEFAULT now(),
                                                  UNIQUE(guild_id, nickname_id)
);

CREATE TABLE IF NOT EXISTS aa_server_nicknames (
    server_id UUID REFERENCES aa_servers(id) ON DELETE CASCADE,
    nickname_id UUID REFERENCES aa_nicknames(id) ON DELETE CASCADE,
    name VARCHAR(18) NOT NULL,
    UNIQUE (server_id, name)
);

CREATE TABLE IF NOT EXISTS tenant_users (
    tenant_id UUID REFERENCES tenants(id) ON DELETE CASCADE,
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    member_at TIMESTAMPTZ DEFAULT now(),
    UNIQUE (tenant_id, user_id)
)

