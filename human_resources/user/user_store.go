package user

import (
	"context"
	"errors"
)

type UserStore interface {
	Clear()
	Get(ctx context.Context, userId UserID) (User, error)
	Save(ctx context.Context, user User) error
	GetMany(ctx context.Context) ([]User, error)
	GetByEmail(ctx context.Context, email string) (User, error)
	// Deprecated: Prefer GetByEmail
	EmailExists(ctx context.Context, email string) (bool, error)
}

var ErrUserNotFound = errors.New("user not found")
