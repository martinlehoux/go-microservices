package group

import (
	"context"
	"go-microservices/common"
	"go-microservices/human_resources/user"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type SqlGroupStore struct {
	conn *pgxpool.Pool
}

func NewSqlGroupStore() SqlGroupStore {
	conn, err := pgxpool.Connect(context.Background(), "postgres://user:password@localhost:5432/human_resources")
	common.PanicOnError(err)
	return SqlGroupStore{
		conn: conn,
	}
}

func (store *SqlGroupStore) Get(ctx context.Context, groupId GroupID) (Group, error) {
	var group Group
	err := store.conn.QueryRow(ctx, "SELECT id, name, description FROM groups WHERE id = $1", groupId).Scan(&group.id, &group.name, &group.description)
	if err == pgx.ErrNoRows {
		return group, ErrGroupNotFound
	}
	members := make([]Membership, 0)
	rows, err := store.conn.Query(ctx, "SELECT user_id, joined_at FROM groups_memberships WHERE group_id = $1", groupId)
	if err != nil {
		return group, err
	}
	defer rows.Close()
	for rows.Next() {
		var membership Membership
		err := rows.Scan(&membership.userID, &membership.joinedAt)
		if err != nil {
			return group, err
		}
		members = append(members, membership)
	}
	group.members = members

	return group, err
}

func (store *SqlGroupStore) Save(ctx context.Context, group Group) error {
	transaction, err := store.conn.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	_, err = transaction.Exec(ctx, "INSERT INTO groups (id, name, description) VALUES ($1, $2, $3) ON CONFLICT (id) DO UPDATE SET name = EXCLUDED.name, description = EXCLUDED.description", group.id, group.name, group.description)
	if err != nil {
		return err
	}
	for _, membership := range group.members {
		_, err := transaction.Exec(ctx, "INSERT INTO groups_memberships (group_id, user_id, joined_at) VALUES ($1, $2, $3) ON CONFLICT (group_id, user_id) DO UPDATE SET joined_at = EXCLUDED.joined_at", group.id, membership.userID, membership.joinedAt)
		if err != nil {
			return err
		}
	}
	err = transaction.Commit(ctx)
	return err
}

func (store *SqlGroupStore) FindByMemberUserId(ctx context.Context, userId user.UserID) ([]GroupDto, error) {
	groups := make([]GroupDto, 0)
	rows, err := store.conn.Query(ctx, "SELECT id, name, description, count(*) as members_count FROM groups JOIN groups_memberships ON groups.id = groups_memberships.group_id WHERE id IN (SELECT group_id FROM groups_memberships WHERE user_id = $1) GROUP BY groups.id, groups.name, groups.description", userId)
	if err != nil {
		return groups, err
	}
	defer rows.Close()
	for rows.Next() {
		var group GroupDto
		err := rows.Scan(&group.ID, &group.Name, &group.Description, &group.MembersCount)
		if err != nil {
			return groups, err
		}
		groups = append(groups, group)
	}
	return groups, err
}

func (store *SqlGroupStore) Clear() {
	_, err := store.conn.Exec(context.Background(), "DELETE FROM groups_memberships")
	common.PanicOnError(err)
	_, err = store.conn.Exec(context.Background(), "DELETE FROM groups")
	common.PanicOnError(err)
}
