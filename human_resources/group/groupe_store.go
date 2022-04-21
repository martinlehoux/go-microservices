package group

import "go-microservices/human_resources/user"

type GroupStore interface {
	Get(groupId GroupID) (Group, error)
	Save(group Group) error
	FindForUser(userId user.UserID) ([]GroupDto, error)
}
