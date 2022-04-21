//go:build spec

package group

import (
	"go-microservices/common"
	"go-microservices/human_resources/user"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddMember(t *testing.T) {
	assert := assert.New(t)

	t.Run("it should add a member to the group", func(t *testing.T) {
		userID := user.UserID{common.CreateID()}
		group := New("test", "test")

		err := group.AddMember(userID)

		assert.NoError(err)
		assert.Equal(1, len(group.members))
	})

	t.Run("it should not add a member to the group if the user is already a member", func(t *testing.T) {
		userID := user.UserID{common.CreateID()}
		group := New("test", "test")

		group.AddMember(userID)
		err := group.AddMember(userID)

		assert.ErrorContains(err, "user already a member")
		assert.Equal(1, len(group.members))
	})
}
