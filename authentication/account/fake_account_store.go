package account

import (
	"context"
	"errors"
)

type FakeAccountStore struct {
	accounts map[AccountID]Account
}

func NewFakeAccountStore() FakeAccountStore {
	return FakeAccountStore{
		accounts: make(map[AccountID]Account),
	}
}

func (store *FakeAccountStore) Save(ctx context.Context, account Account) error {
	store.accounts[account.id] = account
	return nil
}

func (store *FakeAccountStore) LoadForIdentifier(ctx context.Context, identifier string) (Account, error) {
	for _, account := range store.accounts {
		if account.identifier == identifier {
			return account, nil
		}
	}
	return Account{}, errors.New("account not found")
}

func (store *FakeAccountStore) Clear() {
	store.accounts = make(map[AccountID]Account)
}
