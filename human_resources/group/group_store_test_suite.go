//go:build intg

package group

import (
	"go-microservices/common"
	"go-microservices/human_resources/user"
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestPiece struct {
	title string
	run   func(t *testing.T, groupStore TestableGroupStore)
}

type TestableGroupStore interface {
	GroupStore
	// clear the groupStore content
	clear()
}

func GroupStoreTestSuite(t *testing.T, groupStore TestableGroupStore) {
	tests := []TestPiece{
		{title: "SaveAndGet", run: TestSaveAndGet},
		{title: "FindForUser", run: TestFindForUser},
	}
	for _, test := range tests {
		t.Run(test.title, func(t *testing.T) { test.run(t, groupStore) })
	}
}

func TestSaveAndGet(t *testing.T, groupStore TestableGroupStore) {
	assert := assert.New(t)

	t.Run("it should save the Group", func(t *testing.T) {
		t.Cleanup(groupStore.clear)

		group := New("Name", "Description")

		err := groupStore.Save(group)

		assert.NoError(err, "the save should succeed")

		savedGroup, err := groupStore.Get(group.id)
		assert.NoError(err)
		assert.Equal(group, savedGroup)
	})

	t.Run("it should save the latest version of the Group", func(t *testing.T) {
		t.Cleanup(groupStore.clear)

		group := New("Name", "Description")
		groupStore.Save(group)
		group.Rename("New Name")

		err := groupStore.Save(group)

		assert.NoError(err, "the save should succeed")
		savedGroup, _ := groupStore.Get(group.id)
		assert.Equal(group, savedGroup, "the group should be saved exactly")
	})

	t.Run("it should return an error if the Group does not exist", func(t *testing.T) {
		t.Cleanup(groupStore.clear)

		_, err := groupStore.Get(GroupID{common.CreateID()})

		assert.ErrorContains(err, "group not found")
	})

	t.Run("it should save and get group members", func(t *testing.T) {
		t.Cleanup(groupStore.clear)

		group := New("Name", "Description")
		group.AddMember(user.UserID{common.CreateID()})
		group.AddMember(user.UserID{common.CreateID()})

		err := groupStore.Save(group)

		assert.NoError(err, "the save should succeed")

		savedGroup, err := groupStore.Get(group.id)
		assert.NoError(err)
		assert.Equal(2, len(savedGroup.members))
		assert.Equal(group.members[0].userID, savedGroup.members[0].userID)
	})
}

func TestFindForUser(t *testing.T, groupStore TestableGroupStore) {
	assert := assert.New(t)

	t.Run("it should return the groups for the user", func(t *testing.T) {
		t.Cleanup(groupStore.clear)

		userID := user.UserID{common.CreateID()}
		group1 := New("Group 1", "Description")
		group1.AddMember(userID)
		groupStore.Save(group1)
		group2 := New("Group 2", "Description")
		group2.AddMember(userID)
		group2.AddMember(user.UserID{common.CreateID()})
		groupStore.Save(group2)

		groups, err := groupStore.FindForUser(userID)
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
