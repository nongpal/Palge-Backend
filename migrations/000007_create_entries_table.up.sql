CREATE TABLE IF NOT EXISTS entries (
    id bigserial PRIMARY KEY,
    account_id bigint NOT NULL REFERENCES accounts(id),
    amount bigint NOT NULL,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_entries_account_id ON entries (account_id);
