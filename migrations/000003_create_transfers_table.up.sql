CREATE TABLE IF NOT EXISTS transfers (
    id bigserial PRIMARY KEY,
    user_id bigint NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    sender_account_id bigint NOT NULL REFERENCES accounts(id) ON DELETE RESTRICT,
    receiver_account_id bigint NOT NULL REFERENCES accounts(id) ON DELETE RESTRICT,

    amount bigint NOT NULL CHECK (amount > 0),
    idempotency_key text,

    created_at timestamptz NOT NULL DEFAULT NOW(),

    CHECK (sender_account_id <> receiver_account_id)
);

CREATE UNIQUE INDEX IF NOT EXISTS transfer_idempotency_key_idx
    ON transfers(user_id, idempotency_key)
    WHERE idempotency_key IS NOT NULL;

CREATE INDEX IF NOT EXISTS transfers_sender_account_id_idx
    ON transfers(sender_account_id);

CREATE INDEX IF NOT EXISTS transfers_receiver_account_id_idx
    ON transfers(receiver_account_id);

CREATE INDEX IF NOT EXISTS transfer_created_at_idx
    ON transfers(created_at);
