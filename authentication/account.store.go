package authentication

import "go-microservices/common"

type AccountStore interface {
	save(account Account) error
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
