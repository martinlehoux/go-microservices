package group

import (
	"errors"
	"go-microservices/human_resources/user"
)

var ErrGroupNotFound = errors.New("group not found")

type GroupStore interface {
	Get(groupId GroupID) (Group, error)
	Save(group Group) error
	FindForUser(userId user.UserID) ([]GroupDto, error)
}
