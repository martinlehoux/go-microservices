//go:build spec

package authentication

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"go-microservices/authentication/account"
	"go-microservices/common"
	"testing"

	"github.com/stretchr/testify/assert"
)

var ctx = context.Background()

func testModule() *AuthenticationService {
	accountStore := account.NewFakeAccountStore()
	logger := common.NewLogrusLogger()
	privateKey, _ := rsa.GenerateKey(rand.Reader, 1024)
	service := NewAuthenticationService(&accountStore, &logger, *privateKey)
	return &service
}

func TestRegister(t *testing.T) {
	assert := assert.New(t)
	service := testModule()
	t.Cleanup(func() {
		service = testModule()
	})

	t.Run("it should register an account", func(t *testing.T) {
		err := service.Register(ctx, "identifier", "password")

		assert.NoError(err, "the registration should succeed")
		_, err = service.accountStore.GetByIdentifier(ctx, "identifier")
		assert.NoError(err, "the account should be saved")
	})

	t.Run("it should abort if the account already exists", func(t *testing.T) {
		service.accountStore.Save(ctx, account.NewAccount("identifier", []byte("password")))

		err := service.Register(ctx, "identifier", "password")

		assert.ErrorIs(err, ErrIdentifierUsed)
	})
}

func TestAuthenticate(t *testing.T) {
	assert := assert.New(t)
	service := testModule()
	t.Cleanup(func() {
		service = testModule()
	})

	t.Run("it should abort if identifier does not exist", func(t *testing.T) {
		signature, err := service.Authenticate(ctx, "identifier", "password")

		assert.Nil(signature, "the signature should be empty")
		assert.ErrorIs(err, account.ErrAccountNotFound)
	})

	t.Run("it should abort if the password does not match", func(t *testing.T) {
		service.Register(ctx, "identifier", "password")

		signature, err := service.Authenticate(ctx, "identifier", "wrong password")

		assert.Nil(signature)
		assert.ErrorIs(err, ErrWrongPassword)
	})

	t.Run("it should authenticate and return an encrypted token", func(t *testing.T) {
		service.Register(ctx, "identifier", "password")

		signature, err := service.Authenticate(ctx, "identifier", "password")

		assert.NoError(err, "the authentication should succeed")
		assert.NotEmpty(signature, "the signature should exist")
	})
}
