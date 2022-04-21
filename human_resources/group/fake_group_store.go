package group

import (
	"errors"
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
		return group, errors.New("Group not found")
	}
	return group, nil
}

func (store *FakeGroupStore) Save(group Group) error {
	store.groups[group.id] = group
	return nil
}

func (store *FakeGroupStore) FindForUser(userId user.UserID) ([]Group, error) {
	groups := make([]Group, 0)
	for _, group := range store.groups {
		if group.IsMember(userId) {
			groups = append(groups, group)
		}
	}
	return groups, nil
}
