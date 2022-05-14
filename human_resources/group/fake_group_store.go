package group

import (
	"go-microservices/human_resources/user"
)

type FakeGroupStore struct {
	groups map[GroupID]Group
}

func NewFakeGroupStore() FakeGroupStore {
	return FakeGroupStore{
		groups: make(map[GroupID]Group),
	}
}

func (store *FakeGroupStore) Get(groupId GroupID) (Group, error) {
	group, found := store.groups[groupId]
	if !found {
		return group, ErrGroupNotFound
	}
	return group, nil
}

func (store *FakeGroupStore) Save(group Group) error {
	store.groups[group.id] = group
	return nil
}

func (store *FakeGroupStore) FindForUser(userId user.UserID) ([]GroupDto, error) {
	groups := make([]GroupDto, 0)
	for _, group := range store.groups {
		if group.IsMember(userId) {
			groups = append(groups, GroupDto{ID: group.id.String(), Name: group.name, Description: group.description, MembersCount: len(group.members)})
		}
	}
	return groups, nil
}

func (store *FakeGroupStore) Cleanup() {
	store.groups = make(map[GroupID]Group)
}
