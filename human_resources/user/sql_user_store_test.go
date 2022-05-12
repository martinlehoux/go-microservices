//go:build intg

package user

import (
	"context"
	"testing"
)

func (store *SqlUserStore) clear() {
	store.conn.Exec(context.Background(), "DELETE FROM users")
}

func TestSqlUserStore(t *testing.T) {
	store := NewSqlUserStore()
	UserStoreTestSuite(t, &store)
}
