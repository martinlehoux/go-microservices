//go:build spec

package human_resources

import (
	"go-microservices/human_resources/user"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRegister(t *testing.T) {
	assert := assert.New(t)
	userStore := user.NewFakeUserStore()
	service := HumanResourcesService{userStore: &userStore}
	t.Cleanup(func() {
		userStore.Cleanup()
	})

	t.Run("it should successfully register a non existing email", func(t *testing.T) {
		err := service.Register("john@doe.com")

		assert.NoError(err, "expected Register to not error")
		isEmailUsed, _ := userStore.EmailExists("john@doe.com")
		assert.True(isEmailUsed, "expected user to be saved")
	})

	t.Run("it should fail to register an existing email", func(t *testing.T) {
		service.Register("john@doe.com")
		err := service.Register("john@doe.com")

		assert.ErrorContains(err, "email already used", "expected Register to error")
	})
}
