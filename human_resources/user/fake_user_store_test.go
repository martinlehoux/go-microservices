//go:build intg

package user

import "testing"

func (store *FakeUserStore) clear() {
	store.Cleanup()
}

func TestFakeUserStore(t *testing.T) {
	store := NewFakeUserStore()
	UserStoreTestSuite(t, &store)
}
