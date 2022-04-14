package user

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSave(t *testing.T) {
	assert := assert.New(t)
	repository := NewSqlUserStore()
	t.Cleanup(func() {
		repository.truncate()
	})

	t.Run("it should save the User", func(t *testing.T) {
		user := NewUser(NewUserPayload{PreferredName: "joe", Email: "joe@doe.com"})

		err := repository.Save(user)

		assert.Nil(err, "the save should succeed")
	})

	t.Run("it should save the latest version of the User", func(t *testing.T) {
		user := NewUser(NewUserPayload{PreferredName: "paul", Email: "paul@doe.com"})
		repository.Save(user)
		user.Rename("jean paul")

		err := repository.Save(user)

		assert.Nil(err, "the save should succeed")
		savedUser, _ := repository.Load(user.Id)
		assert.Equal(user, savedUser, "the user should be saved exactly")
	})
}

func TestEmailExists(t *testing.T) {
	assert := assert.New(t)
	repository := NewSqlUserStore()
	t.Cleanup(func() {
		repository.truncate()
	})

	t.Run("it should return true if the email exists", func(t *testing.T) {
		repository.Save(NewUser(NewUserPayload{PreferredName: "john", Email: "john@doe.com"}))

		exists, err := repository.EmailExists("john@doe.com")

		assert.Nil(err, "the query should not fail")
		assert.True(exists, "the email should exist")
	})

	t.Run("it should return false if the email does not exist", func(t *testing.T) {
		exists, err := repository.EmailExists("john@king.com")

		assert.Nil(err, "the query should not fail")
		assert.False(exists, "the email should not exist")
	})
}

func TestGetMany(t *testing.T) {
	assert := assert.New(t)
	repository := NewSqlUserStore()
	t.Cleanup(func() {
		repository.truncate()
	})

	t.Run("it should return all the Users", func(t *testing.T) {
		repository.Save(NewUser(NewUserPayload{PreferredName: "john", Email: "john@travolta.com"}))
		repository.Save(NewUser(NewUserPayload{PreferredName: "jane", Email: "jane@roosevelt.com"}))

		users, err := repository.GetMany()

		assert.Nil(err, "the get should succeed")
		assert.Equal(2, len(users), "there should be two users")
		assert.Equal("john", users[0].PreferredName, "the first user should be john")
		assert.Equal("jane", users[1].PreferredName, "the second user should be jane")
	})
}
