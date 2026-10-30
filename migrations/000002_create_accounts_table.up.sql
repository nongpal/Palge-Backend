CREATE TABLE IF NOT EXISTS accounts (
  id bigserial PRIMARY KEY,
  user_id bigint NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  balance bigint NOT NULL CHECK ( balance >= 0 )
);

CREATE INDEX idx_accounts_user_id ON accounts (user_id);
