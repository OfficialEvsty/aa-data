ALTER TABLE lunark DROP COLUMN opened;
ALTER TABLE tenant_lunark ADD COLUMN opened BOOLEAN DEFAULT TRUE;

CREATE UNIQUE INDEX uniq_open_lunark_per_tenant
    ON tenant_lunark (tenant_id)
    WHERE opened = true;