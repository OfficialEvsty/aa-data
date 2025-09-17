ALTER TABLE salaries ADD COLUMN is_deleted BOOLEAN NOT NULL DEFAULT FALSE;

CREATE TABLE IF NOT EXISTS tenant_chains (
                                             tenant_id UUID REFERENCES tenants(id) ON DELETE CASCADE,
                                             chain_id UUID REFERENCES chains(chain_id) ON DELETE CASCADE,
                                             UNIQUE(tenant_id, chain_id)
)