CREATE TABLE IF NOT EXISTS transfers (
    id bigserial PRIMARY KEY,
    sender_account_id bigint NOT NULL REFERENCES accounts(id),
    receiver_account_id bigint NOT NULL REFERENCES accounts(id),
    amount bigint NOT NULL CHECK (amount > 0),
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_transfers_sender_account ON transfers(sender_account_id);
CREATE INDEX idx_transfers_receiver_account ON transfers(receiver_account_id);
