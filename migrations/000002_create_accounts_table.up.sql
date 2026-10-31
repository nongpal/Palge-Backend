CREATE TABLE IF NOT EXISTS accounts (
  id bigserial PRIMARY KEY,
  user_id bigint NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  balance bigint NOT NULL CHECK ( balance >= 0 ),
  status text NOT NULL DEFAULT "ACTIVATE",
  created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
  updated_at timestamp(0) with time zone NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_accounts_user_id ON accounts (user_id);
