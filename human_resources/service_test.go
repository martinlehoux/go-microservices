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
		err := service.Register("john@doe.com", "John")

		assert.NoError(err, "expected Register to not error")
		user, _ := userStore.GetByEmail("john@doe.com")
		assert.Equal(user.GetEmail(), "john@doe.com", "the email should match")
	})

	t.Run("it should fail to register an existing email", func(t *testing.T) {
		service.Register("john@doe.com", "John")
		err := service.Register("john@doe.com", "John")

		assert.ErrorContains(err, "email already used", "expected Register to error")
	})
}
