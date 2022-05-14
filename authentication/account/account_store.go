package account

import "context"

type AccountStore interface {
	Clear()
	Save(ctx context.Context, account Account) error
	LoadForIdentifier(ctx context.Context, identifier string) (Account, error)
}
