package user

import (
	"errors"
)

type UserStore interface {
	load(userId UserID) (User, error)
	save(user User) error
	getMany() ([]User, error)
	emailExists(email string) (bool, error)
}

type FakeUserStore struct {
	users map[UserID]User
}

func bootstrapFakeUserStore() FakeUserStore {
	return FakeUserStore{
		users: make(map[UserID]User),
	}
}

func (store *FakeUserStore) load(userId UserID) (User, error) {
	user, found := store.users[userId]
	if !found {
		return user, errors.New("User not found")
	}
	return user, nil
}

func (store *FakeUserStore) save(user User) error {
	store.users[user.Id] = user
	return nil
}

func (store *FakeUserStore) getMany() ([]User, error) {
	users := make([]User, 0)
	for _, user := range store.users {
		users = append(users, user)
	}
	return users, nil
}

func (store *FakeUserStore) emailExists(email string) (bool, error) {
	for _, user := range store.users {
		if user.Email == email {
			return true, nil
		}
	}
	return false, nil
}
