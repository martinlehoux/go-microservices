package authentication

import (
	"errors"
	"go-microservices/common"
)

type AccountStore interface {
	save(account Account) error
	loadForIdentifier(identifier string) (Account, error)
}

type FakeAccountStore struct {
	accounts map[common.ID]Account
}

func bootstrapFakeAccountStore() FakeAccountStore {
	return FakeAccountStore{
		accounts: make(map[common.ID]Account),
	}
}

func (store *FakeAccountStore) save(account Account) error {
	store.accounts[account.Id] = account
	return nil
}

func (store *FakeAccountStore) loadForIdentifier(identifier string) (Account, error) {
	for _, account := range store.accounts {
		if account.Identifier == identifier {
			return account, nil
		}
	}
	return Account{}, errors.New("account not found")
}
