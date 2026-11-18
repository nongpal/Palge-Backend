CREATE TABLE IF NOT EXISTS entries (
    id bigserial PRIMARY KEY,
    account_id bigint NOT NULL REFERENCES accounts(id),
    transfer_id bigint NOT NULL REFERENCES transfers(id) ON DELETE RESTRICT,
    amount bigint NOT NULL,
    balance_after bigint NOT NULL DEFAULT 0,
    kind text NOT NULL DEFAULT 'TRANSFER_IN' CHECK (kind IN ('TRANSFER_IN', 'TRANSFER_OUT', 'DEPOSIT', 'WITHDRAW')),
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_entries_account_id ON entries (account_id);
