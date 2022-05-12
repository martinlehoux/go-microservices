//go:build intg

package group

import (
	"context"
	"testing"
)

func (store *SqlGroupStore) clear() {
	store.conn.Exec(context.Background(), "DELETE FROM groups")
	store.conn.Exec(context.Background(), "DELETE FROM group_memberships")
}

func TestSqlGroupStore(t *testing.T) {
	store := NewSqlGroupStore()
	GroupStoreTestSuite(t, &store)
}
