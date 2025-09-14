INSERT INTO aa_bosses (id, name, drop, img_url) VALUES (1111, '', '[]'::jsonb, '') ON CONFLICT (id) DO NOTHING;

INSERT INTO aa_template_events (id, name) VALUES (14, 'Ашьяра, Гленн и Лорея');

INSERT INTO aa_event_bosses (event_template_id, boss_id) VALUES (14, 1111) ON CONFLICT (event_template_id, boss_id) DO NOTHING;

INSERT INTO aa_servers (id, name, external_id) VALUES ('1e0c8de5-9fc0-4538-b433-b9eae8e3a776'::UUID, 'server', 15);

INSERT INTO aa_nicknames (id, name, server_id) VALUES ('b6d3afdf-1a9d-4bc4-a8b6-d9377889c25b'::UUID, 'finick', '1e0c8de5-9fc0-4538-b433-b9eae8e3a776'::UUID), ('1d3a19ba-4dfe-47fc-8e3d-bfe7d63c72fa'::UUID, 'drugfinika', '1e0c8de5-9fc0-4538-b433-b9eae8e3a776'::UUID);

INSERT INTO aa_guilds (id, name, server_id) VALUES ('5f25e254-d493-4779-8903-0e67aae7443c'::UUID, 'fanem', '1e0c8de5-9fc0-4538-b433-b9eae8e3a776'::UUID);

INSERT INTO aa_guild_nicknames (guild_id, nickname_id) VALUES ('5f25e254-d493-4779-8903-0e67aae7443c'::UUID, 'b6d3afdf-1a9d-4bc4-a8b6-d9377889c25b'::UUID), ('5f25e254-d493-4779-8903-0e67aae7443c'::UUID, '1d3a19ba-4dfe-47fc-8e3d-bfe7d63c72fa'::UUID)