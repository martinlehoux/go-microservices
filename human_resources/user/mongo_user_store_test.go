//go:build intg

package user

import (
	"testing"
)

func TestMongoUserStore(t *testing.T) {
	store := NewMongoUserStore()
	UserStoreTestSuite(t, &store)
}
