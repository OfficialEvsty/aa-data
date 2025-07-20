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

CREATE TABLE IF NOT EXISTS aa_guild_events (
                                               guild_id UUID REFERENCES aa_guilds(id) ON DELETE CASCADE,
                                               event_id BIGSERIAL REFERENCES events(id) ON DELETE CASCADE,
                                               UNIQUE (guild_id, event_id)
);

CREATE TABLE IF NOT EXISTS aa_guild_storage (
                                                guild_id UUID REFERENCES aa_guilds(id) ON DELETE CASCADE,
                                                storage_id UUID REFERENCES storage(id) ON DELETE CASCADE,
                                                UNIQUE(guild_id, storage_id)
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

