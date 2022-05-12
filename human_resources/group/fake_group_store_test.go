//go:build intg

package group

import "testing"

func (store *FakeGroupStore) clear() {
	store.groups = make(map[GroupID]Group)
}

func TestFakeGroupStore(t *testing.T) {
	store := NewFakeGroupStore()
	GroupStoreTestSuite(t, &store)
}
