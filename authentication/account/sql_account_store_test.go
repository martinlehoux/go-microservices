package account

import (
	"go-microservices/common"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSave(t *testing.T) {
	assert := assert.New(t)
	repository := NewSqlAccountStore()
	t.Cleanup(func() {
		repository.truncate()
	})

	t.Run("it should save the Account", func(t *testing.T) {
		identifier := common.CreateID().String()
		account := NewAccount(identifier, []byte("password"))

		err := repository.Save(account)

		assert.Nil(err, "the save should succeed")
		savedAccount, err := repository.LoadForIdentifier(identifier)
		assert.Nil(err, "the account should be saved")
		assert.Equal(account, savedAccount, "the account should be saved exactly")
	})

	t.Run("it should save the latest version of the Account", func(t *testing.T) {
		t.Skip("there is not update for now")
		identifier := common.CreateID().String()
		account := NewAccount(identifier, []byte("password"))
		repository.Save(account)
	})
}

func TestLoadForIdentifier(t *testing.T) {
	assert := assert.New(t)
	repository := NewSqlAccountStore()
	t.Cleanup(func() {
		repository.truncate()
	})

	t.Run("it should not get an Account with antoher identifier", func(t *testing.T) {
		identifier := common.CreateID().String()
		repository.Save(NewAccount(identifier, []byte("password")))

		account, err := repository.LoadForIdentifier("wrong")

		assert.Equal(account, Account{}, "the account should be empty")
		assert.Error(err, "not found")
	})

	t.Run("it should get an Account with the correct identifier", func(t *testing.T) {
		identifier := common.CreateID().String()
		savedAccount := NewAccount(identifier, []byte("password"))
		repository.Save(savedAccount)

		account, err := repository.LoadForIdentifier(identifier)

		assert.Nil(err, "the load should succeed")
		assert.Equal(account, savedAccount, "the account should be empty")
	})

}
