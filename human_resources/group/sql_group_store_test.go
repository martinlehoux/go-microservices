//go:build intg

package group

import (
	"testing"
)

func TestSqlGroupStore(t *testing.T) {
	store := NewSqlGroupStore("test_human_resources")
	GroupStoreTestSuite(t, &store)
}
