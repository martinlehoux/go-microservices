//go:build intg

package user

import (
	"testing"
)

func TestSqlUserStore(t *testing.T) {
	store := NewSqlUserStore("test_human_resources")
	UserStoreTestSuite(t, &store)
}
