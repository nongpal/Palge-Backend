CREATE EXTENSION IF NOT EXISTS citext;

CREATE TABLE IF NOT EXISTS users (
  id bigserial PRIMARY KEY,
  created_at timestamptz NOT NULL DEFAULT NOW(),
  updated_at timestamptz NOT NULL DEFAULT NOW(),

  name text NOT NULL CHECK (length(btrim(name)) > 0),
  email citext NOT NULL UNIQUE CHECK (email <> ''),
  password_hash bytea NOT NULL,

  activated boolean NOT NULL DEFAULT false,
  version integer NOT NULL DEFAULT 1 CHECK (version > 0)
);
