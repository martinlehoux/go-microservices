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

func TestSave(t *testing.T) {
	assert := assert.New(t)
	repository := NewSqlUserStore()

	t.Run("it should save the User", func(t *testing.T) {
		t.Cleanup(repository.truncate)

		user := NewUser(NewUserPayload{PreferredName: "joe", Email: "joe@doe.com"})

		err := repository.Save(user)

		assert.NoError(err, "the save should succeed")
	})

	t.Run("it should save the latest version of the User", func(t *testing.T) {
		t.Cleanup(repository.truncate)

		user := NewUser(NewUserPayload{PreferredName: "paul", Email: "paul@doe.com"})
		repository.Save(user)
		user.Rename("jean paul")

		err := repository.Save(user)

		assert.NoError(err, "the save should succeed")
		savedUser, _ := repository.Load(user.id)
		assert.Equal(user, savedUser, "the user should be saved exactly")
	})
}

func TestEmailExists(t *testing.T) {
	assert := assert.New(t)
	repository := NewSqlUserStore()

	t.Run("it should return true if the email exists", func(t *testing.T) {
		t.Cleanup(repository.truncate)

		repository.Save(NewUser(NewUserPayload{PreferredName: "john", Email: "john@doe.com"}))

		exists, err := repository.EmailExists("john@doe.com")

		assert.NoError(err, "the query should not fail")
		assert.True(exists, "the email should exist")
	})

	t.Run("it should return false if the email does not exist", func(t *testing.T) {
		t.Cleanup(repository.truncate)

		exists, err := repository.EmailExists("john@king.com")

		assert.NoError(err, "the query should not fail")
		assert.False(exists, "the email should not exist")
	})
}

func TestGetMany(t *testing.T) {
	assert := assert.New(t)
	repository := NewSqlUserStore()

	t.Run("it should return all the Users", func(t *testing.T) {
		t.Cleanup(repository.truncate)

		repository.Save(NewUser(NewUserPayload{PreferredName: "john", Email: "john@travolta.com"}))
		repository.Save(NewUser(NewUserPayload{PreferredName: "jane", Email: "jane@roosevelt.com"}))

		users, err := repository.GetMany()

		assert.NoError(err, "the get should succeed")
		assert.Equal(2, len(users), "there should be two users")
		assert.Equal("john", users[0].preferredName, "the first user should be john")
		assert.Equal("jane", users[1].preferredName, "the second user should be jane")
	})

	t.Run("it should return an empty slice if there are no users", func(t *testing.T) {
		t.Cleanup(repository.truncate)

		users, err := repository.GetMany()

		assert.NoError(err, "the get should succeed")
		assert.Equal(make([]User, 0), users, "there should be no users")
	})
}

func TestGetByEmail(t *testing.T) {
	assert := assert.New(t)
	repository := NewSqlUserStore()

	t.Run("it should return the User", func(t *testing.T) {
		t.Cleanup(repository.truncate)

		repository.Save(NewUser(NewUserPayload{PreferredName: "john", Email: "john@doe.com"}))

		user, err := repository.GetByEmail("john@doe.com")

		assert.NoError(err, "the get should succeed")
		assert.Equal("john", user.preferredName, "the user should be john")
	})
}
