CREATE EXTENSION IF NOT EXISTS citext;

CREATE TABLE IF NOT EXISTS users (
  id bigserial PRIMARY KEY,
  created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
  name text NOT NULL,
  email CITEXT UNIQUE NOT NUll,
  password_hash bytea NOT NULL,
  -- rencana default ke false
  -- activated bool NOT NULL DEFAULT false,
  activated bool NOT NULL,
  version integer NOT NULL DEFAULT 1
);
