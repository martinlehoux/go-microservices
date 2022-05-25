package user

import (
	"context"
	"errors"
)

type UserStore interface {
	Clear()
	Get(ctx context.Context, userId UserID) (User, error)
	Save(ctx context.Context, user User) error
	GetMany(ctx context.Context) ([]UserDto, error)
	GetByEmail(ctx context.Context, email string) (User, error)
}

var ErrUserNotFound = errors.New("user not found")
