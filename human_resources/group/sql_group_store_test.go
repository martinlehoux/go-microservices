//go:build intg

package group

import (
	"testing"
)

func TestSqlGroupStore(t *testing.T) {
	store := NewSqlGroupStore()
	GroupStoreTestSuite(t, &store)
}
