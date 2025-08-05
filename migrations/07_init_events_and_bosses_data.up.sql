-- PRE-Initialize bosses ids
-- Гленн
INSERT INTO aa_bosses (id, name, drop, img_url) VALUES (13279, '', '[]'::jsonb, '') ON CONFLICT (id) DO NOTHING;
-- Лорея
INSERT INTO aa_bosses (id, name, drop, img_url) VALUES (13278, '', '[]'::jsonb, '') ON CONFLICT (id) DO NOTHING;
-- Ашьяра
INSERT INTO aa_bosses (id, name, drop, img_url) VALUES (12649, '', '[]'::jsonb, '') ON CONFLICT (id) DO NOTHING;
-- Т2 Гленн
INSERT INTO aa_bosses (id, name, drop, img_url) VALUES (20352, '', '[]'::jsonb, '') ON CONFLICT (id) DO NOTHING;
-- Т2 Лорея
INSERT INTO aa_bosses (id, name, drop, img_url) VALUES (20353, '', '[]'::jsonb, '') ON CONFLICT (id) DO NOTHING;
-- Т2 Ашьяра
INSERT INTO aa_bosses (id, name, drop, img_url) VALUES (20357, '', '[]'::jsonb, '') ON CONFLICT (id) DO NOTHING;
-- Морфеус
INSERT INTO aa_bosses (id, name, drop, img_url) VALUES (8766, '', '[]'::jsonb, '') ON CONFLICT (id) DO NOTHING;
-- Марли
INSERT INTO aa_bosses (id, name, drop, img_url) VALUES (8554, '', '[]'::jsonb, '') ON CONFLICT (id) DO NOTHING;
-- Разъяренная Сехекмет
INSERT INTO aa_bosses (id, name, drop, img_url) VALUES (22069, '', '[]'::jsonb, '') ON CONFLICT (id) DO NOTHING;
-- Голиаф, механический скарабей
INSERT INTO aa_bosses (id, name, drop, img_url) VALUES (20449, '', '[]'::jsonb, '') ON CONFLICT (id) DO NOTHING;
-- Кракен
INSERT INTO aa_bosses (id, name, drop, img_url) VALUES (7607, '', '[]'::jsonb, '') ON CONFLICT (id) DO NOTHING;
-- Авиара
INSERT INTO aa_bosses (id, name, drop, img_url) VALUES (18230, '', '[]'::jsonb, '') ON CONFLICT (id) DO NOTHING;
-- Каллеиль
INSERT INTO aa_bosses (id, name, drop, img_url) VALUES (20343, '', '[]'::jsonb, '') ON CONFLICT (id) DO NOTHING;
-- Ксанатос
INSERT INTO aa_bosses (id, name, drop, img_url) VALUES (18653, '', '[]'::jsonb, '') ON CONFLICT (id) DO NOTHING;
-- Левиафан
INSERT INTO aa_bosses (id, name, drop, img_url) VALUES (14915, '', '[]'::jsonb, '') ON CONFLICT (id) DO NOTHING;
-- Калидис
INSERT INTO aa_bosses (id, name, drop, img_url) VALUES (19790, '', '[]'::jsonb, '') ON CONFLICT (id) DO NOTHING;
-- Анталлон, призрак войны
INSERT INTO aa_bosses (id, name, drop, img_url) VALUES (19871, '', '[]'::jsonb, '') ON CONFLICT (id) DO NOTHING;



-- Initialize template events
INSERT INTO aa_template_events (id, name) VALUES (0, 'Ашьяра, Гленн и Лорея');
INSERT INTO aa_template_events (id, name) VALUES (1, 'Т2 Ашьяра');
INSERT INTO aa_template_events (id, name) VALUES (2, 'Т2 Гленн и Лорея');
INSERT INTO aa_template_events (id, name) VALUES (3, 'Морфеус');
INSERT INTO aa_template_events (id, name) VALUES (4, 'Марли');
INSERT INTO aa_template_events (id, name) VALUES (5, 'Сехекмет');
INSERT INTO aa_template_events (id, name) VALUES (6, 'Голиаф');
INSERT INTO aa_template_events (id, name) VALUES (7, 'Кракен');
INSERT INTO aa_template_events (id, name) VALUES (8, 'Авиара');
INSERT INTO aa_template_events (id, name) VALUES (9, 'Каллеиль');
INSERT INTO aa_template_events (id, name) VALUES (10, 'Ксанатос');
INSERT INTO aa_template_events (id, name) VALUES (11, 'Левиафан');
INSERT INTO aa_template_events (id, name) VALUES (12, 'Калидис');
INSERT INTO aa_template_events (id, name) VALUES (13, 'Анталлон, призрак войны');

-- Initialize template event's bosses
-- Ашьяра (ID: 12649), Гленн (ID: 13279) и Лорея (ID: 13278)
INSERT INTO aa_event_bosses (event_template_id, boss_id) VALUES (0, 13279) ON CONFLICT (event_template_id, boss_id) DO NOTHING;
INSERT INTO aa_event_bosses (event_template_id, boss_id) VALUES (0, 13278) ON CONFLICT (event_template_id, boss_id) DO NOTHING;
INSERT INTO aa_event_bosses (event_template_id, boss_id) VALUES (0, 12649) ON CONFLICT (event_template_id, boss_id) DO NOTHING;
-- Т2 Ашьяра (ID: 20357)
INSERT INTO aa_event_bosses (event_template_id, boss_id) VALUES (1, 20357) ON CONFLICT (event_template_id, boss_id) DO NOTHING;
-- T2 Гленн (ID: 20352) и Лорея (ID: 20353)
INSERT INTO aa_event_bosses (event_template_id, boss_id) VALUES (2, 20353) ON CONFLICT (event_template_id, boss_id) DO NOTHING;
INSERT INTO aa_event_bosses (event_template_id, boss_id) VALUES (2, 20352) ON CONFLICT (event_template_id, boss_id) DO NOTHING;
-- Морфеус (ID: 8766)
INSERT INTO aa_event_bosses (event_template_id, boss_id) VALUES (3, 8766) ON CONFLICT (event_template_id, boss_id) DO NOTHING;
-- Марли (ID: 8554)
INSERT INTO aa_event_bosses (event_template_id, boss_id) VALUES (4, 8554) ON CONFLICT (event_template_id, boss_id) DO NOTHING;
-- Т2 Сехекмет (ID: 22069)
INSERT INTO aa_event_bosses (event_template_id, boss_id) VALUES (5, 22069) ON CONFLICT (event_template_id, boss_id) DO NOTHING;
-- Голиаф (ID: 20449)
INSERT INTO aa_event_bosses (event_template_id, boss_id) VALUES (6, 20449) ON CONFLICT (event_template_id, boss_id) DO NOTHING;
-- Кракен (ID: 7607)
INSERT INTO aa_event_bosses (event_template_id, boss_id) VALUES (7, 7607) ON CONFLICT (event_template_id, boss_id) DO NOTHING;
-- Авиара (ID: 18230)
INSERT INTO aa_event_bosses (event_template_id, boss_id) VALUES (8, 18230) ON CONFLICT (event_template_id, boss_id) DO NOTHING;
-- Каллеиль (ID: 20343)
INSERT INTO aa_event_bosses (event_template_id, boss_id) VALUES (9, 20343) ON CONFLICT (event_template_id, boss_id) DO NOTHING;
-- Ксанатос (ID: 18653)
INSERT INTO aa_event_bosses (event_template_id, boss_id) VALUES (10, 18653) ON CONFLICT (event_template_id, boss_id) DO NOTHING;
-- Левиафан (ID: 14915)
INSERT INTO aa_event_bosses (event_template_id, boss_id) VALUES (11, 14915) ON CONFLICT (event_template_id, boss_id) DO NOTHING;
-- Калидис (ID: 19790)
INSERT INTO aa_event_bosses (event_template_id, boss_id) VALUES (12, 19790) ON CONFLICT (event_template_id, boss_id) DO NOTHING;
-- Анталлон, призрак войны ID: 19871)
INSERT INTO aa_event_bosses (event_template_id, boss_id) VALUES (13, 19871) ON CONFLICT (event_template_id, boss_id) DO NOTHING;