BEGIN;

CREATE TABLE users (
  id UUID PRIMARY KEY NOT NULL,
  preferred_name TEXT NOT NULL,
  email TEXT NOT NULL
);

CREATE UNIQUE INDEX users_email_unique ON users (email);

COMMIT;