//go:build spec

package authentication

import (
	"crypto/rand"
	"crypto/rsa"
	"go-microservices/authentication/account"
	"go-microservices/common"
	"testing"

	"github.com/stretchr/testify/assert"
)

func testModule() *AuthenticationService {
	accountStore := account.NewFakeAccountStore()
	privateKey, _ := rsa.GenerateKey(rand.Reader, 1024)
	service := AuthenticationService{accountStore: &accountStore, privateKey: *privateKey}
	return &service
}

func TestRegister(t *testing.T) {
	assert := assert.New(t)
	service := testModule()
	t.Cleanup(func() {
		service = testModule()
	})

	t.Run("it should register an account", func(t *testing.T) {
		err := service.Register("identifier", "password")

		assert.Nil(err, "the registration should succeed")
		_, err = service.accountStore.LoadForIdentifier("identifier")
		assert.Nil(err, "the account should be saved")
	})

	t.Run("it should abort if the account already exists", func(t *testing.T) {
		service.accountStore.Save(account.NewAccount("identifier", []byte("password")))

		err := service.Register("identifier", "password")

		assert.ErrorContains(err, "identifier already used", "the registration should fail")
	})
}

func TestAuthenticate(t *testing.T) {
	assert := assert.New(t)
	service := testModule()
	t.Cleanup(func() {
		service = testModule()
	})

	t.Run("it should abort if identifier does not exist", func(t *testing.T) {
		token, signature, err := service.Authenticate("identifier", "password")

		assert.Nil(signature)
		assert.Equal(token, common.Token{})
		assert.ErrorContains(err, "account not found", "the authentication should fail")
	})

	t.Run("it should abort if the password does not match", func(t *testing.T) {
		service.Register("identifier", "password")

		token, signature, err := service.Authenticate("identifier", "wrong password")

		assert.Nil(signature)
		assert.Equal(token, common.Token{})
		assert.ErrorContains(err, "password mismatch", "the authentication should fail")
	})

	t.Run("it should authenticate and return an encrypted token", func(t *testing.T) {
		service.Register("identifier", "password")

		token, signature, err := service.Authenticate("identifier", "password")

		assert.Nil(err, "the authentication should succeed")
		assert.Equal(token.Identifier, "identifier", "the token should be encrypted")
		assert.NotEmpty(signature, "the signature should exist")
	})
}
