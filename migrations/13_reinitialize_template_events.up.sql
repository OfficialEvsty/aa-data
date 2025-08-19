-- Remove wrong event-bosses relation
-- Т2 Ашьяра (ID: 20357)
DELETE FROM aa_event_bosses WHERE boss_id=20357;
-- T2 Гленн (ID: 20352) и Лорея (ID: 20353)
DELETE FROM aa_event_bosses WHERE boss_id=20353;
DELETE FROM aa_event_bosses WHERE boss_id=20352;
-- Т2 Ашьяра (ID: 20357)
INSERT INTO aa_event_bosses (event_template_id, boss_id) VALUES (0, 20357) ON CONFLICT (event_template_id, boss_id) DO NOTHING;
-- T2 Гленн (ID: 20352) и Лорея (ID: 20353)
INSERT INTO aa_event_bosses (event_template_id, boss_id) VALUES (0, 20353) ON CONFLICT (event_template_id, boss_id) DO NOTHING;
INSERT INTO aa_event_bosses (event_template_id, boss_id) VALUES (0, 20352) ON CONFLICT (event_template_id, boss_id) DO NOTHING;
-- Remove wrong templates
DELETE FROM aa_template_events WHERE id = 1 OR id = 2
-- RE-Initialize template events
-- INSERT INTO aa_template_events (id, name) VALUES (0, 'Ашьяра, Гленн и Лорея');
-- INSERT INTO aa_template_events (id, name) VALUES (3, 'Морфеус');
-- INSERT INTO aa_template_events (id, name) VALUES (4, 'Марли');
-- INSERT INTO aa_template_events (id, name) VALUES (5, 'Сехекмет');
-- INSERT INTO aa_template_events (id, name) VALUES (6, 'Голиаф');
-- INSERT INTO aa_template_events (id, name) VALUES (7, 'Кракен');
-- INSERT INTO aa_template_events (id, name) VALUES (8, 'Авиара');
-- INSERT INTO aa_template_events (id, name) VALUES (9, 'Каллеиль');
-- INSERT INTO aa_template_events (id, name) VALUES (10, 'Ксанатос');
-- INSERT INTO aa_template_events (id, name) VALUES (11, 'Левиафан');
-- INSERT INTO aa_template_events (id, name) VALUES (12, 'Калидис');
-- INSERT INTO aa_template_events (id, name) VALUES (13, 'Анталлон, призрак войны');