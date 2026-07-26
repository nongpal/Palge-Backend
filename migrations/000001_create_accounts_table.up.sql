CREATE TABLE IF NOT EXISTS accounts (
  id bigserial PRIMARY KEY,
  owner text NOT NULL,
  balance bigint NOT NULL CHECK ( balance >= 0 )
);


