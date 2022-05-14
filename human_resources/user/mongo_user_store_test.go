//go:build intg

package user

import (
	"context"
	"testing"

	"go.mongodb.org/mongo-driver/bson"
)

func (store *MongoUserStore) clear() {
	store.collection.DeleteMany(context.Background(), bson.D{})
}

func TestMongoUserStore(t *testing.T) {
	store := NewMongoUserStore()
	UserStoreTestSuite(t, &store)
}
