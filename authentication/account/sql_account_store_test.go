//go:build intg

package account

import (
	"context"
	"go-microservices/common"
	"testing"

	"github.com/stretchr/testify/assert"
)

var ctx = context.Background()

func TestSave(t *testing.T) {
	assert := assert.New(t)
	store := NewSqlAccountStore()

	t.Run("it should save the Account", func(t *testing.T) {
		t.Cleanup(store.Clear)

		identifier := common.CreateID().String()
		account := NewAccount(identifier, []byte("password"))

		err := store.Save(ctx, account)

		assert.NoError(err, "the save should succeed")
		savedAccount, err := store.LoadForIdentifier(ctx, identifier)
		assert.NoError(err, "the account should be saved")
		assert.Equal(account, savedAccount, "the account should be saved exactly")
	})

	t.Run("it should save the latest version of the Account", func(t *testing.T) {
		t.Skip("there is not update for now")
		t.Cleanup(store.Clear)

		identifier := common.CreateID().String()
		account := NewAccount(identifier, []byte("password"))
		store.Save(ctx, account)
	})
}

func TestLoadForIdentifier(t *testing.T) {
	assert := assert.New(t)
	store := NewSqlAccountStore()

	t.Run("it should not get an Account with antoher identifier", func(t *testing.T) {
		t.Cleanup(store.Clear)

		identifier := common.CreateID().String()
		store.Save(ctx, NewAccount(identifier, []byte("password")))

		account, err := store.LoadForIdentifier(ctx, "wrong")

		assert.Equal(account, Account{}, "the account should be empty")
		assert.ErrorContains(err, "no rows in result set", "the error should be returned")
	})

	t.Run("it should get an Account with the correct identifier", func(t *testing.T) {
		t.Cleanup(store.Clear)

		identifier := common.CreateID().String()
		savedAccount := NewAccount(identifier, []byte("password"))
		store.Save(ctx, savedAccount)

		account, err := store.LoadForIdentifier(ctx, identifier)

		assert.NoError(err, "the load should succeed")
		assert.Equal(account, savedAccount, "the account should be empty")
	})

}
