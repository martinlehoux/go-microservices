package group

import (
	"context"
	"errors"
	"go-microservices/human_resources/user"
)

var ErrGroupNotFound = errors.New("group not found")

type GroupStore interface {
	Clear()
	Get(ctx context.Context, groupId GroupID) (Group, error)
	Save(ctx context.Context, group Group) error
	FindByMemberUserId(ctx context.Context, userId user.UserID) ([]GroupDto, error)
}
