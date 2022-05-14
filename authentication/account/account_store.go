package account

import (
	"context"
	"errors"
)

type AccountStore interface {
	Clear()
	Save(ctx context.Context, account Account) error
	GetByIdentifier(ctx context.Context, identifier string) (Account, error)
}

var ErrAccountNotFound = errors.New("account not found")
