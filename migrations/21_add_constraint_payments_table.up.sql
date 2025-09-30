ALTER TABLE payments
    ADD CONSTRAINT unique_salary_chain_id UNIQUE (salary_id, root_chain_id);