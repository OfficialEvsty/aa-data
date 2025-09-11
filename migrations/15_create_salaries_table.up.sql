CREATE TABLE IF NOT EXISTS chains (
    chain_id UUID PRIMARY KEY,
    parent_chain_id UUID REFERENCES chains(chain_id) ON DELETE CASCADE,
    nickname_id UUID REFERENCES aa_nicknames(id) ON DELETE CASCADE,
    chained_at TIMESTAMPTZ DEFAULT now(),
    active BOOLEAN NOT NULL DEFAULT TRUE,
    UNIQUE(chain_id, nickname_id)
);

CREATE TABLE IF NOT EXISTS salaries (
    id UUID PRIMARY KEY,
    fond BIGINT DEFAULT 0,
    min_attendance INTEGER NOT NULL DEFAULT 1,
    tax INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ DEFAULT now()
);

CREATE TABLE IF NOT EXISTS payments (
    salary_id UUID REFERENCES salaries(id) ON DELETE CASCADE,
    root_chain_id UUID REFERENCES chains(chain_id) ON DELETE CASCADE,
    salary BIGINT NOT NULL,
    reason TEXT DEFAULT ''
);

CREATE TABLE IF NOT EXISTS lunark_salaries (
    lunark_id UUID REFERENCES lunark(id) ON DELETE CASCADE,
    salary_id UUID REFERENCES salaries(id) ON DELETE CASCADE,
    paid_at TIMESTAMPTZ DEFAULT now(),
    UNIQUE(lunark_id, salary_id)
);

CREATE TABLE IF NOT EXISTS excluded_participants_salary (
    root_chain_id UUID REFERENCES chains(chain_id) ON DELETE CASCADE,
    salary_id UUID REFERENCES salaries(id) ON DELETE CASCADE,
    reason TEXT DEFAULT '',
    UNIQUE(root_chain_id, salary_id)
);