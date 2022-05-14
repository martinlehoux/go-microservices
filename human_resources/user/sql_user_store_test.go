//go:build intg

package user

import (
	"testing"
)

func TestSqlUserStore(t *testing.T) {
	store := NewSqlUserStore()
	UserStoreTestSuite(t, &store)
}
