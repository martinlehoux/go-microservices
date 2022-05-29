//go:build intg

package group

import (
	"testing"
)

func TestMongoGroupStore(t *testing.T) {
	store := NewMongoGroupStore("test_human_resources")
	GroupStoreTestSuite(t, &store)
}
