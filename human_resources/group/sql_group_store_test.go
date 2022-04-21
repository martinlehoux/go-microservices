//go:build intg

package group

import (
	"context"
	"go-microservices/common"
	"go-microservices/human_resources/user"
	"testing"

	"github.com/stretchr/testify/assert"
)

func (store *SqlGroupStore) truncate() {
	store.conn.Exec(context.Background(), "DELETE FROM groups")
	store.conn.Exec(context.Background(), "DELETE FROM group_memberships")
}

func TestSaveAndGet(t *testing.T) {
	assert := assert.New(t)
	store := NewSqlGroupStore()

	t.Run("it should save the Group", func(t *testing.T) {
		t.Cleanup(store.truncate)

		group := New("Name", "Description")

		err := store.Save(group)

		assert.NoError(err, "the save should succeed")

		savedGroup, err := store.Get(group.id)
		assert.NoError(err)
		assert.Equal(group, savedGroup)
	})

	t.Run("it should save the latest version of the Group", func(t *testing.T) {
		t.Cleanup(store.truncate)

		group := New("Name", "Description")
		store.Save(group)
		group.Rename("New Name")

		err := store.Save(group)

		assert.NoError(err, "the save should succeed")
		savedGroup, _ := store.Get(group.id)
		assert.Equal(group, savedGroup, "the group should be saved exactly")
	})

	t.Run("it should return an error if the Group does not exist", func(t *testing.T) {
		t.Cleanup(store.truncate)

		_, err := store.Get(GroupID{common.CreateID()})

		assert.ErrorContains(err, "group not found")
	})

	t.Run("it should save and get group members", func(t *testing.T) {
		t.Cleanup(store.truncate)

		group := New("Name", "Description")
		group.AddMember(user.UserID{common.CreateID()})
		group.AddMember(user.UserID{common.CreateID()})

		err := store.Save(group)

		assert.NoError(err, "the save should succeed")

		savedGroup, err := store.Get(group.id)
		assert.NoError(err)
		assert.Equal(2, len(savedGroup.members))
		assert.Equal(group.members[0].userID, savedGroup.members[0].userID)
	})
}

func TestFindForUser(t *testing.T) {
	assert := assert.New(t)
	store := NewSqlGroupStore()

	t.Run("it should return the groups for the user", func(t *testing.T) {
		t.Cleanup(store.truncate)

		userID := user.UserID{common.CreateID()}
		group1 := New("Group 1", "Description")
		group1.AddMember(userID)
		store.Save(group1)
		group2 := New("Group 2", "Description")
		group2.AddMember(userID)
		group2.AddMember(user.UserID{common.CreateID()})
		store.Save(group2)

		groups, err := store.FindForUser(userID)
		assert.NoError(err)
		assert.Equal(2, len(groups))
		assert.Equal(GroupDto{
			ID:           group1.id.String(),
			Name:         group1.name,
			Description:  group1.description,
			MembersCount: 1,
		}, groups[0])
		assert.Equal(GroupDto{
			ID:           group2.id.String(),
			Name:         group2.name,
			Description:  group2.description,
			MembersCount: 2,
		}, groups[1])
	})
}
