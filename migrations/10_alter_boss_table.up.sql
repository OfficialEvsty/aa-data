ALTER TABLE aa_bosses DROP COLUMN img_grade_url;

CREATE TABLE IF NOT EXISTS user_activities (
                                               raid_id UUID REFERENCES raids(id) ON DELETE CASCADE,
                                               user_id UUID REFERENCES users(id) ON DELETE CASCADE,
                                               UNIQUE (raid_id, user_id)
);