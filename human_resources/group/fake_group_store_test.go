//go:build intg

package group

import "testing"

func TestFakeGroupStore(t *testing.T) {
	store := NewFakeGroupStore()
	GroupStoreTestSuite(t, &store)
}
