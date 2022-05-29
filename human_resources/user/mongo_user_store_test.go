//go:build intg

package user

import (
	"testing"
)

func TestMongoUserStore(t *testing.T) {
	store := NewMongoUserStore("test_human_resources")
	UserStoreTestSuite(t, &store)
}
