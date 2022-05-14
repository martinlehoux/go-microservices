//go:build intg

package user

import "testing"

func TestFakeUserStore(t *testing.T) {
	store := NewFakeUserStore()
	UserStoreTestSuite(t, &store)
}
