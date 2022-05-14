//go:build intg

package group

import (
	"context"
	"go-microservices/common"
	"go-microservices/human_resources/user"
	"testing"

	"github.com/stretchr/testify/assert"
)

var ctx = context.Background()

type TestPiece struct {
	title string
	run   func(t *testing.T, groupStore GroupStore)
}

func GroupStoreTestSuite(t *testing.T, groupStore GroupStore) {
	tests := []TestPiece{
		{title: "SaveAndGet", run: TestSaveAndGet},
		{title: "FindByMemberUserId", run: TestFindByMemberUserId},
	}
	for _, test := range tests {
		t.Run(test.title, func(t *testing.T) { test.run(t, groupStore) })
	}
}

func TestSaveAndGet(t *testing.T, groupStore GroupStore) {
	assert := assert.New(t)

	t.Run("it should save the Group", func(t *testing.T) {
		t.Cleanup(groupStore.Clear)

		group := New("Name", "Description")

		err := groupStore.Save(ctx, group)

		assert.NoError(err, "the save should succeed")

		savedGroup, err := groupStore.Get(ctx, group.id)
		assert.NoError(err)
		assert.Equal(group, savedGroup)
	})

	t.Run("it should save the latest version of the Group", func(t *testing.T) {
		t.Cleanup(groupStore.Clear)

		group := New("Name", "Description")
		groupStore.Save(ctx, group)
		group.Rename("New Name")

		err := groupStore.Save(ctx, group)

		assert.NoError(err, "the save should succeed")
		savedGroup, _ := groupStore.Get(ctx, group.id)
		assert.Equal(group, savedGroup, "the group should be saved exactly")
	})

	t.Run("it should return an error if the Group does not exist", func(t *testing.T) {
		t.Cleanup(groupStore.Clear)

		_, err := groupStore.Get(ctx, GroupID{common.CreateID()})

		assert.ErrorIs(err, ErrGroupNotFound)
	})

	t.Run("it should save and get group members", func(t *testing.T) {
		t.Cleanup(groupStore.Clear)

		group := New("Name", "Description")
		group.AddMember(user.UserID{common.CreateID()})
		group.AddMember(user.UserID{common.CreateID()})

		err := groupStore.Save(ctx, group)

		assert.NoError(err, "the save should succeed")

		savedGroup, err := groupStore.Get(ctx, group.id)
		assert.NoError(err)
		assert.Equal(2, len(savedGroup.members))
		assert.Equal(group.members[0].userID, savedGroup.members[0].userID)
	})
}

func TestFindByMemberUserId(t *testing.T, groupStore GroupStore) {
	assert := assert.New(t)

	t.Run("it should return the groups for the user", func(t *testing.T) {
		t.Cleanup(groupStore.Clear)

		userID := user.UserID{common.CreateID()}
		group1 := New("Group 1", "Description")
		group1.AddMember(userID)
		groupStore.Save(ctx, group1)
		group2 := New("Group 2", "Description")
		group2.AddMember(userID)
		group2.AddMember(user.UserID{common.CreateID()})
		groupStore.Save(ctx, group2)

		groups, err := groupStore.FindByMemberUserId(ctx, userID)
		assert.NoError(err)
		assert.Equal(2, len(groups))
		assert.ElementsMatch([]GroupDto{{
			ID:           group1.id.String(),
			Name:         group1.name,
			Description:  group1.description,
			MembersCount: 1,
		}, {
			ID:           group2.id.String(),
			Name:         group2.name,
			Description:  group2.description,
			MembersCount: 2,
		}}, groups)
	})
}
