BEGIN;

CREATE TABLE groups (
  id UUID PRIMARY KEY NOT NULL,
  name TEXT NOT NULL,
  description TEXT NOT NULL
);

CREATE TABLE groups_memberships (
  group_id UUID NOT NULL,
  user_id UUID NOT NULL,
  joined_at TIMESTAMP NOT NULL,
  FOREIGN KEY (group_id) REFERENCES groups (id)
);

CREATE UNIQUE INDEX groups_memberships_group_id_user_id_key ON groups_memberships (group_id, user_id);

COMMIT;