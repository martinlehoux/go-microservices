CREATE TABLE accounts (
  id UUID NOT NULL PRIMARY KEY,
  identifier TEXT NOT NULL,
  hashed_password BYTEA NOT NULL
);