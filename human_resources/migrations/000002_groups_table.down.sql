BEGIN;

DROP INDEX groups_memberships_group_id_user_id_key;

DROP TABLE groups_memberships;

DROP TABLE groups;

COMMIT;