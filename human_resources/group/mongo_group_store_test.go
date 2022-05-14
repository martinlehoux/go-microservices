//go:build intg

package group

import (
	"context"
	"testing"

	"go.mongodb.org/mongo-driver/bson"
)

func (store *MongoGroupStore) clear() {
	store.collection.DeleteMany(context.Background(), bson.D{})
}

func TestMongoGroupStore(t *testing.T) {
	store := NewMongoGroupStore()
	GroupStoreTestSuite(t, &store)
}
