CREATE TABLE IF NOT EXISTS accounts (
  id bigserial PRIMARY KEY,
  user_id bigint NOT NULL REFERENCES users(id) ON DELETE CASCADE,

  balance bigint NOT NULL DEFAULT 0 CHECK (balance >= 0),
  status text NOT NULL DEFAULT 'ACTIVE'
    CHECK (status IN ('ACTIVE', 'FROZEN', 'CLOSED')),

  created_at timestamptz NOT NULL DEFAULT NOW(),
  updated_at timestamptz NOT NULL DEFAULT NOW(),

  version integer NOT NULL DEFAULT 1 CHECK (version > 0)
);

CREATE INDEX idx_accounts_user_id ON accounts (user_id);
