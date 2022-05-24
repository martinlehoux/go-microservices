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
	logger := common.NewLogrusLogger()
	userStore := user.NewFakeUserStore()
	groupStore := group.NewFakeGroupStore()
	service := NewHumanResourcesService(&userStore, &groupStore, &logger)

	t.Run("it should successfully register a non existing email", func(t *testing.T) {
		t.Cleanup(func() {
			userStore.Clear()
			groupStore.Clear()
		})
		err := service.Register(ctx, "john@doe.com", "John")

		assert.NoError(err, "expected Register to not error")
		user, _ := userStore.GetByEmail(ctx, "john@doe.com")
		assert.Equal(user.GetEmail(), "john@doe.com", "the email should match")
	})

	t.Run("it should fail to register an existing email", func(t *testing.T) {
		t.Cleanup(func() {
			userStore.Clear()
			groupStore.Clear()
		})
		service.Register(ctx, "john@doe.com", "John")
		err := service.Register(ctx, "john@doe.com", "John")

		assert.ErrorIs(err, ErrEmailUsed, "expected Register to error")
	})
}

func TestUserJoinGroup(t *testing.T) {
	assert := assert.New(t)
	logger := common.NewLogrusLogger()
	userStore := user.NewFakeUserStore()
	groupStore := group.NewFakeGroupStore()
	service := NewHumanResourcesService(&userStore, &groupStore, &logger)

	t.Run("it should fail to join a group if the group does not exist", func(t *testing.T) {
		t.Cleanup(func() {
			userStore.Clear()
			groupStore.Clear()
		})
		user := user.New(user.NewUserPayload{Email: "test@test.com", PreferredName: "Test"})

		err := service.UserJoinGroup(ctx, user.GetID(), group.GroupID{common.CreateID()})

		assert.ErrorIs(err, group.ErrGroupNotFound, "expected UserJoinGroup to error")
	})

	t.Run("it should fail to join a group if the user does not exist", func(t *testing.T) {
		t.Cleanup(func() {
			userStore.Clear()
			groupStore.Clear()
		})
		groupToJoin := group.New("Group", "")
		groupStore.Save(ctx, groupToJoin)

		err := service.UserJoinGroup(ctx, user.UserID{common.CreateID()}, groupToJoin.GetID())

		assert.ErrorIs(err, user.ErrUserNotFound, "expected UserJoinGroup to error")
	})

	t.Run("it should make the user join the group", func(t *testing.T) {
		t.Cleanup(func() {
			userStore.Clear()
			groupStore.Clear()
		})
		groupToJoin := group.New("Group", "")
		userToJoin := user.New(user.NewUserPayload{Email: "test@test.com", PreferredName: "Test"})
		userStore.Save(ctx, userToJoin)
		groupStore.Save(ctx, groupToJoin)

		err := service.UserJoinGroup(ctx, userToJoin.GetID(), groupToJoin.GetID())

		assert.NoError(err, "expected UserJoinGroup to succeed")
		updatedGroup, _ := groupStore.Get(ctx, groupToJoin.GetID())
		assert.True(updatedGroup.IsMember(userToJoin.GetID()), "expected the user to be a member of the group")
	})
}
