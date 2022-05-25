//go:build intg

package user

import (
	"context"
	"go-microservices/common"
	"testing"

	"github.com/stretchr/testify/assert"
)

var ctx = context.Background()

type TestPiece struct {
	title string
	run   func(t *testing.T, userStore UserStore)
}

func UserStoreTestSuite(t *testing.T, userStore UserStore) {
	tests := []TestPiece{
		{title: "SaveAndGet", run: TestSaveAndGet},
		{title: "GetMany", run: TestGetMany},
		{title: "GetByEmail", run: TestGetByEmail},
	}
	for _, test := range tests {
		t.Run(test.title, func(t *testing.T) { test.run(t, userStore) })
	}
}

func TestSaveAndGet(t *testing.T, userStore UserStore) {
	assert := assert.New(t)

	t.Run("it should save the User", func(t *testing.T) {
		t.Cleanup(userStore.Clear)

		user := New("joe", "joe@doe.com")

		err := userStore.Save(ctx, user)

		assert.NoError(err, "the save should succeed")

		user, err = userStore.Get(ctx, user.id)
		assert.NoError(err, "the get should succeed")
		assert.Equal("joe", user.preferredName, "the user should be joe")
	})

	t.Run("it should save the latest version of the User", func(t *testing.T) {
		t.Cleanup(userStore.Clear)

		user := New("paul", "paul@doe.com")
		userStore.Save(ctx, user)
		user.Rename("jean paul")

		err := userStore.Save(ctx, user)

		assert.NoError(err, "the save should succeed")
		savedUser, _ := userStore.Get(ctx, user.id)
		assert.Equal(user, savedUser, "the user should be saved exactly")
	})

	t.Run("it should return an error if the User is not found", func(t *testing.T) {
		t.Cleanup(userStore.Clear)

		user, err := userStore.Get(ctx, UserID{common.CreateID()})

		assert.ErrorIs(ErrUserNotFound, err, "the get should fail")
		assert.Equal(User{}, user, "the user should be empty")
	})
}

func TestGetMany(t *testing.T, userStore UserStore) {
	assert := assert.New(t)

	t.Run("it should return all the Users", func(t *testing.T) {
		t.Cleanup(userStore.Clear)

		john := New("john", "john@travolta.com")
		jane := New("jane", "jane@roosevelt.com")
		userStore.Save(ctx, john)
		userStore.Save(ctx, jane)

		users, err := userStore.GetMany(ctx)

		assert.NoError(err)
		assert.Equal(2, len(users))
		assert.ElementsMatch([]User{john, jane}, users)
	})

	t.Run("it should return an empty slice if there are no users", func(t *testing.T) {
		t.Cleanup(userStore.Clear)

		users, err := userStore.GetMany(ctx)

		assert.NoError(err, "the get should succeed")
		assert.Equal(make([]User, 0), users, "there should be no users")
	})
}

func TestGetByEmail(t *testing.T, userStore UserStore) {
	assert := assert.New(t)

	t.Run("it should return the User", func(t *testing.T) {
		t.Cleanup(userStore.Clear)

		userStore.Save(ctx, New("john", "john@doe.com"))

		user, err := userStore.GetByEmail(ctx, "john@doe.com")

		assert.NoError(err, "the get should succeed")
		assert.Equal("john", user.preferredName, "the user should be john")
	})

	t.Run("it should return an error if the User is not found", func(t *testing.T) {
		t.Cleanup(userStore.Clear)

		user, err := userStore.GetByEmail(ctx, "john@doe.com")

		assert.ErrorIs(err, ErrUserNotFound)
		assert.Equal(User{}, user, "the user should be empty")
	})
}
