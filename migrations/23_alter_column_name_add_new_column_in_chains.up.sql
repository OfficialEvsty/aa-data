ALTER TABLE chains RENAME COLUMN chained_at TO created_at;
ALTER TABLE chains ADD COLUMN chained_at TIMESTAMPTZ NULL;