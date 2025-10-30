ALTER TABLE requests ADD COLUMN source_id UUID NULL;
ALTER TABLE requests ADD COLUMN source_name TEXT DEFAULT 'system' NOT NULL;