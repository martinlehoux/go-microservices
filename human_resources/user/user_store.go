package user

import "errors"

type UserStore interface {
	Clear()
	Get(userId UserID) (User, error)
	Save(user User) error
	GetMany() ([]User, error)
	GetByEmail(email string) (User, error)
	// Deprecated: Prefer GetByEmail
	EmailExists(email string) (bool, error)
}

var ErrUserNotFound = errors.New("user not found")
