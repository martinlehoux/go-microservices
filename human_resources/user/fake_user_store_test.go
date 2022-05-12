//go:build intg

package user

import "testing"

func (store *FakeUserStore) clear() {
	store.users = make(map[UserID]User)
}

func TestFakeUserStore(t *testing.T) {
	store := NewFakeUserStore()
	UserStoreTestSuite(t, &store)
}
