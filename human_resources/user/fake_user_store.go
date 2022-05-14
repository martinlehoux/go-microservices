package user

import "errors"

type FakeUserStore struct {
	users map[UserID]User
}

func NewFakeUserStore() FakeUserStore {
	return FakeUserStore{
		users: make(map[UserID]User),
	}
}

func (store *FakeUserStore) Get(userId UserID) (User, error) {
	user, found := store.users[userId]
	if !found {
		return user, ErrUserNotFound
	}
	return user, nil
}

func (store *FakeUserStore) Save(user User) error {
	store.users[user.id] = user
	return nil
}

func (store *FakeUserStore) GetMany() ([]User, error) {
	users := make([]User, 0)
	for _, user := range store.users {
		users = append(users, user)
	}
	return users, nil
}

func (store *FakeUserStore) EmailExists(email string) (bool, error) {
	for _, user := range store.users {
		if user.email == email {
			return true, nil
		}
	}
	return false, nil
}

func (store *FakeUserStore) GetByEmail(email string) (User, error) {
	for _, user := range store.users {
		if user.email == email {
			return user, nil
		}
	}
	return User{}, errors.New("user not found")
}

func (store *FakeUserStore) Cleanup() {
	store.users = make(map[UserID]User)
}
