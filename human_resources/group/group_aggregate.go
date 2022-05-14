package group

import (
	"errors"
	"go-microservices/common"
	"go-microservices/human_resources/user"
	"time"
)

type GroupID struct{ common.ID }

type Membership struct {
	userID   user.UserID
	joinedAt time.Time
}

type Group struct {
	id          GroupID
	name        string
	description string
	members     []Membership
}

var ErrAlreadyMember = errors.New("user already a member")

func New(name string, description string) Group {
	return Group{
		id:          GroupID{common.CreateID()},
		name:        name,
		description: description,
		members:     make([]Membership, 0),
	}
}

func (group *Group) AddMember(userID user.UserID) error {
	if group.IsMember(userID) {
		return ErrAlreadyMember
	}
	group.members = append(group.members, Membership{
		userID:   userID,
		joinedAt: time.Now(),
	})
	return nil
}

func (group Group) IsMember(userID user.UserID) bool {
	for _, membership := range group.members {
		if membership.userID == userID {
			return true
		}
	}
	return false
}

func (group *Group) Rename(name string) {
	group.name = name
}

func (group Group) GetID() GroupID {
	return group.id
}
