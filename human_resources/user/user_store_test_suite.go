//go:build intg

package user

import (
	"go-microservices/common"
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestPiece struct {
	title string
	run   func(t *testing.T, userStore TestableUserStore)
}

type TestableUserStore interface {
	UserStore
	// clear the userStore content
	clear()
}

func UserStoreTestSuite(t *testing.T, userStore TestableUserStore) {
	tests := []TestPiece{
		{title: "SaveAndGet", run: TestSaveAndGet},
		{title: "EmailExists", run: TestEmailExists},
		{title: "GetMany", run: TestGetMany},
		{title: "GetByEmail", run: TestGetByEmail},
	}
	for _, test := range tests {
		t.Run(test.title, func(t *testing.T) { test.run(t, userStore) })
	}
}

func TestSaveAndGet(t *testing.T, userStore TestableUserStore) {
	assert := assert.New(t)

	t.Run("it should save the User", func(t *testing.T) {
		t.Cleanup(userStore.clear)

		user := New(NewUserPayload{PreferredName: "joe", Email: "joe@doe.com"})

		err := userStore.Save(user)

		assert.NoError(err, "the save should succeed")

		user, err = userStore.Get(user.id)
		assert.NoError(err, "the get should succeed")
		assert.Equal("joe", user.preferredName, "the user should be joe")
	})

	t.Run("it should save the latest version of the User", func(t *testing.T) {
		t.Cleanup(userStore.clear)

		user := New(NewUserPayload{PreferredName: "paul", Email: "paul@doe.com"})
		userStore.Save(user)
		user.Rename("jean paul")

		err := userStore.Save(user)

		assert.NoError(err, "the save should succeed")
		savedUser, _ := userStore.Get(user.id)
		assert.Equal(user, savedUser, "the user should be saved exactly")
	})

	t.Run("it should return an error if the User is not found", func(t *testing.T) {
		t.Cleanup(userStore.clear)

		user, err := userStore.Get(UserID{common.CreateID()})

		assert.ErrorIs(ErrUserNotFound, err, "the get should fail")
		assert.Equal(User{}, user, "the user should be empty")
	})
}

func TestEmailExists(t *testing.T, userStore TestableUserStore) {
	assert := assert.New(t)

	t.Run("it should return true if the email exists", func(t *testing.T) {
		t.Cleanup(userStore.clear)

		userStore.Save(New(NewUserPayload{PreferredName: "john", Email: "john@doe.com"}))

		exists, err := userStore.EmailExists("john@doe.com")

		assert.NoError(err, "the query should not fail")
		assert.True(exists, "the email should exist")
	})

	t.Run("it should return false if the email does not exist", func(t *testing.T) {
		t.Cleanup(userStore.clear)

		exists, err := userStore.EmailExists("john@king.com")

		assert.NoError(err, "the query should not fail")
		assert.False(exists, "the email should not exist")
	})
}

func TestGetMany(t *testing.T, userStore TestableUserStore) {
	assert := assert.New(t)

	t.Run("it should return all the Users", func(t *testing.T) {
		t.Cleanup(userStore.clear)

		john := New(NewUserPayload{PreferredName: "john", Email: "john@travolta.com"})
		jane := New(NewUserPayload{PreferredName: "jane", Email: "jane@roosevelt.com"})
		userStore.Save(john)
		userStore.Save(jane)

		users, err := userStore.GetMany()

		assert.NoError(err)
		assert.Equal(2, len(users))
		assert.ElementsMatch([]User{john, jane}, users)
	})

	t.Run("it should return an empty slice if there are no users", func(t *testing.T) {
		t.Cleanup(userStore.clear)

		users, err := userStore.GetMany()

		assert.NoError(err, "the get should succeed")
		assert.Equal(make([]User, 0), users, "there should be no users")
	})
}

func TestGetByEmail(t *testing.T, userStore TestableUserStore) {
	assert := assert.New(t)

	t.Run("it should return the User", func(t *testing.T) {
		t.Cleanup(userStore.clear)

		userStore.Save(New(NewUserPayload{PreferredName: "john", Email: "john@doe.com"}))

		user, err := userStore.GetByEmail("john@doe.com")

		assert.NoError(err, "the get should succeed")
		assert.Equal("john", user.preferredName, "the user should be john")
	})
}
