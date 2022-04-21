//go:build intg

package account

import (
	"context"
	"go-microservices/common"
	"testing"

	"github.com/stretchr/testify/assert"
)

func (store *SqlAccountStore) truncate() error {
	_, err := store.conn.Exec(context.Background(), "DELETE FROM accounts")
	return err
}

func TestSave(t *testing.T) {
	assert := assert.New(t)
	store := NewSqlAccountStore()
	t.Cleanup(func() {
		store.truncate()
	})

	t.Run("it should save the Account", func(t *testing.T) {
		identifier := common.CreateID().String()
		account := NewAccount(identifier, []byte("password"))

		err := store.Save(account)

		assert.NoError(err, "the save should succeed")
		savedAccount, err := store.LoadForIdentifier(identifier)
		assert.NoError(err, "the account should be saved")
		assert.Equal(account, savedAccount, "the account should be saved exactly")
	})

	t.Run("it should save the latest version of the Account", func(t *testing.T) {
		t.Skip("there is not update for now")
		identifier := common.CreateID().String()
		account := NewAccount(identifier, []byte("password"))
		store.Save(account)
	})
}

func TestLoadForIdentifier(t *testing.T) {
	assert := assert.New(t)
	store := NewSqlAccountStore()
	t.Cleanup(func() {
		store.truncate()
	})

	t.Run("it should not get an Account with antoher identifier", func(t *testing.T) {
		identifier := common.CreateID().String()
		store.Save(NewAccount(identifier, []byte("password")))

		account, err := store.LoadForIdentifier("wrong")

		assert.Equal(account, Account{}, "the account should be empty")
		assert.ErrorContains(err, "no rows in result set", "the error should be returned")
	})

	t.Run("it should get an Account with the correct identifier", func(t *testing.T) {
		identifier := common.CreateID().String()
		savedAccount := NewAccount(identifier, []byte("password"))
		store.Save(savedAccount)

		account, err := store.LoadForIdentifier(identifier)

		assert.NoError(err, "the load should succeed")
		assert.Equal(account, savedAccount, "the account should be empty")
	})

}
