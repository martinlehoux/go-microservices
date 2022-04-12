package account

import (
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

func (store *FakeAccountStore) Save(account Account) error {
	store.accounts[account.Id] = account
	return nil
}

func (store *FakeAccountStore) LoadForIdentifier(identifier string) (Account, error) {
	for _, account := range store.accounts {
		if account.Identifier == identifier {
			return account, nil
		}
	}
	return Account{}, errors.New("account not found")
}
