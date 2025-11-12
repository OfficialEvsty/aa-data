CREATE TABLE IF NOT EXISTS archived_chains (
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    chain_id UUID NOT NULL REFERENCES chains(chain_id) ON DELETE CASCADE,
    archived_at TIMESTAMPTZ DEFAULT now() NOT NULL,
    UNIQUE(chain_id)
)