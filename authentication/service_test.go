package authentication

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func testModule() *AuthenticationService {
	accountStore := bootstrapFakeAccountStore()
	service := AuthenticationService{accountStore: &accountStore}
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
		_, err = service.accountStore.loadForIdentifier("identifier")
		assert.Nil(err, "the account should be saved")
	})

	t.Run("it should abort if the account already exists", func(t *testing.T) {
		service.accountStore.save(NewAccount("identifier", []byte("password")))

		err := service.Register("identifier", "password")

		assert.ErrorContains(err, "identifier already used", "the registration should fail")
	})
}
