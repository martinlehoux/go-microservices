//go:build intg

package user

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func (store *SqlUserStore) truncate() {
	store.conn.Exec(context.Background(), "DELETE FROM users")
}

func TestSaveAndGet(t *testing.T) {
	assert := assert.New(t)
	store := NewSqlUserStore()

	t.Run("it should save the User", func(t *testing.T) {
		t.Cleanup(store.truncate)

		user := New(NewUserPayload{PreferredName: "joe", Email: "joe@doe.com"})

		err := store.Save(user)

		assert.NoError(err, "the save should succeed")

		user, err = store.Get(user.id)
		assert.NoError(err, "the get should succeed")
		assert.Equal("joe", user.preferredName, "the user should be joe")
	})

	t.Run("it should save the latest version of the User", func(t *testing.T) {
		t.Cleanup(store.truncate)

		user := New(NewUserPayload{PreferredName: "paul", Email: "paul@doe.com"})
		store.Save(user)
		user.Rename("jean paul")

		err := store.Save(user)

		assert.NoError(err, "the save should succeed")
		savedUser, _ := store.Get(user.id)
		assert.Equal(user, savedUser, "the user should be saved exactly")
	})
}

func TestEmailExists(t *testing.T) {
	assert := assert.New(t)
	store := NewSqlUserStore()

	t.Run("it should return true if the email exists", func(t *testing.T) {
		t.Cleanup(store.truncate)

		store.Save(New(NewUserPayload{PreferredName: "john", Email: "john@doe.com"}))

		exists, err := store.EmailExists("john@doe.com")

		assert.NoError(err, "the query should not fail")
		assert.True(exists, "the email should exist")
	})

	t.Run("it should return false if the email does not exist", func(t *testing.T) {
		t.Cleanup(store.truncate)

		exists, err := store.EmailExists("john@king.com")

		assert.NoError(err, "the query should not fail")
		assert.False(exists, "the email should not exist")
	})
}

func TestGetMany(t *testing.T) {
	assert := assert.New(t)
	store := NewSqlUserStore()

	t.Run("it should return all the Users", func(t *testing.T) {
		t.Cleanup(store.truncate)

		store.Save(New(NewUserPayload{PreferredName: "john", Email: "john@travolta.com"}))
		store.Save(New(NewUserPayload{PreferredName: "jane", Email: "jane@roosevelt.com"}))

		users, err := store.GetMany()

		assert.NoError(err, "the get should succeed")
		assert.Equal(2, len(users), "there should be two users")
		assert.Equal("john", users[0].preferredName, "the first user should be john")
		assert.Equal("jane", users[1].preferredName, "the second user should be jane")
	})

	t.Run("it should return an empty slice if there are no users", func(t *testing.T) {
		t.Cleanup(store.truncate)

		users, err := store.GetMany()

		assert.NoError(err, "the get should succeed")
		assert.Equal(make([]User, 0), users, "there should be no users")
	})
}

func TestGetByEmail(t *testing.T) {
	assert := assert.New(t)
	store := NewSqlUserStore()

	t.Run("it should return the User", func(t *testing.T) {
		t.Cleanup(store.truncate)

		store.Save(New(NewUserPayload{PreferredName: "john", Email: "john@doe.com"}))

		user, err := store.GetByEmail("john@doe.com")

		assert.NoError(err, "the get should succeed")
		assert.Equal("john", user.preferredName, "the user should be john")
	})
}
