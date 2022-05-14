//go:build intg

package group

import (
	"testing"
)

func TestMongoGroupStore(t *testing.T) {
	store := NewMongoGroupStore()
	GroupStoreTestSuite(t, &store)
}
