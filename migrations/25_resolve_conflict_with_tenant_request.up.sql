ALTER TABLE tenant_requests DROP COLUMN user_id;
ALTER TABLE tenant_requests ADD COLUMN request_id UUID REFERENCES requests(id) ON DELETE CASCADE;