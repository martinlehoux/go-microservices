package account

import (
	"context"
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

func (store *FakeAccountStore) GetByIdentifier(ctx context.Context, identifier string) (Account, error) {
	for _, account := range store.accounts {
		if account.identifier == identifier {
			return account, nil
		}
	}
	return Account{}, ErrAccountNotFound
}

func (store *FakeAccountStore) Clear() {
	store.accounts = make(map[AccountID]Account)
}
