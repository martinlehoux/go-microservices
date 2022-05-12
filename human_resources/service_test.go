//go:build spec

package human_resources

import (
	"go-microservices/common"
	"go-microservices/human_resources/group"
	"go-microservices/human_resources/user"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRegister(t *testing.T) {
	assert := assert.New(t)
	userStore := user.NewFakeUserStore()
	groupStore := group.NewFakeGroupStore()
	service := NewHumanResourcesService(&userStore, nil)

	t.Run("it should successfully register a non existing email", func(t *testing.T) {
		t.Cleanup(func() {
			userStore.Cleanup()
			groupStore.Cleanup()
		})
		err := service.Register("john@doe.com", "John")

		assert.NoError(err, "expected Register to not error")
		user, _ := userStore.GetByEmail("john@doe.com")
		assert.Equal(user.GetEmail(), "john@doe.com", "the email should match")
	})

	t.Run("it should fail to register an existing email", func(t *testing.T) {
		t.Cleanup(func() {
			userStore.Cleanup()
			groupStore.Cleanup()
		})
		service.Register("john@doe.com", "John")
		err := service.Register("john@doe.com", "John")

		assert.ErrorContains(err, "email already used", "expected Register to error")
	})
}

func TestUserJoinGroup(t *testing.T) {
	assert := assert.New(t)
	userStore := user.NewFakeUserStore()
	groupStore := group.NewFakeGroupStore()
	service := NewHumanResourcesService(&userStore, &groupStore)

	t.Run("it should fail to join a group if the group does not exist", func(t *testing.T) {
		t.Cleanup(func() {
			userStore.Cleanup()
			groupStore.Cleanup()
		})
		user := user.New(user.NewUserPayload{Email: "test@test.com", PreferredName: "Test"})

		err := service.UserJoinGroup(user.GetID(), group.GroupID{common.CreateID()})

		assert.ErrorContains(err, "group not found", "expected UserJoinGroup to error")
	})

	t.Run("it should fail to join a group if the user does not exist", func(t *testing.T) {
		t.Cleanup(func() {
			userStore.Cleanup()
			groupStore.Cleanup()
		})
		groupToJoin := group.New("Group", "")
		groupStore.Save(groupToJoin)

		err := service.UserJoinGroup(user.UserID{common.CreateID()}, groupToJoin.GetID())

		assert.ErrorContains(err, "user not found", "expected UserJoinGroup to error")
	})

	t.Run("it should make the user join the group", func(t *testing.T) {
		t.Cleanup(func() {
			userStore.Cleanup()
			groupStore.Cleanup()
		})
		groupToJoin := group.New("Group", "")
		userToJoin := user.New(user.NewUserPayload{Email: "test@test.com", PreferredName: "Test"})
		userStore.Save(userToJoin)
		groupStore.Save(groupToJoin)

		err := service.UserJoinGroup(userToJoin.GetID(), groupToJoin.GetID())

		assert.NoError(err, "expected UserJoinGroup to succeed")
		updatedGroup, _ := groupStore.Get(groupToJoin.GetID())
		assert.True(updatedGroup.IsMember(userToJoin.GetID()), "expected the user to be a member of the group")
	})
}
