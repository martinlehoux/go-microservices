package user

import (
	"context"
)

type FakeUserStore struct {
	users map[UserID]User
}

func NewFakeUserStore() FakeUserStore {
	return FakeUserStore{
		users: make(map[UserID]User),
	}
}

func (store *FakeUserStore) Get(ctx context.Context, userId UserID) (User, error) {
	user, found := store.users[userId]
	if !found {
		return user, ErrUserNotFound
	}
	return user, nil
}

func (store *FakeUserStore) Save(ctx context.Context, user User) error {
	store.users[user.id] = user
	return nil
}

func (store *FakeUserStore) GetMany(ctx context.Context) ([]UserDto, error) {
	users := make([]UserDto, 0)
	for _, user := range store.users {
		users = append(users, UserDto{
			ID:            user.id.String(),
			PreferredName: user.preferredName,
			Email:         user.email,
		})
	}
	return users, nil
}

func (store *FakeUserStore) GetByEmail(ctx context.Context, email string) (User, error) {
	for _, user := range store.users {
		if user.email == email {
			return user, nil
		}
	}
	return User{}, ErrUserNotFound
}

func (store *FakeUserStore) Clear() {
	store.users = make(map[UserID]User)
}
