package human_resources

import (
	"context"
	"errors"
	"go-microservices/common"
	"go-microservices/human_resources/group"
	"go-microservices/human_resources/user"
)

var (
	ErrEmailUsed = errors.New("email already used")
)

type HumanResourcesService struct {
	userStore  user.UserStore
	groupStore group.GroupStore
	logger     common.Logger
}

func NewHumanResourcesService(userStore user.UserStore, groupStore group.GroupStore, logger common.Logger) HumanResourcesService {
	return HumanResourcesService{userStore: userStore, groupStore: groupStore, logger: logger}
}

func (service *HumanResourcesService) Register(ctx context.Context, email string, preferredName string) error {
	var err error
	service.logger.With(ctx, "email", email)
	service.logger.Info(ctx, "starting user registration for email %s", email)

	_, err = service.userStore.GetByEmail(ctx, email)

	switch err {
	case nil:
		err = ErrEmailUsed
		service.logger.Info(ctx, "failed to register user for email %s: %s", email)
		return err
	case user.ErrUserNotFound:
		break
	default:
		service.logger.Warn(ctx, "failed to register user for email %s: %s", email)
		return err
	}

	user := user.New(user.NewUserPayload{Email: email, PreferredName: preferredName})
	service.logger.With(ctx, "userId", user.GetID().String())

	err = service.userStore.Save(ctx, user)
	if err != nil {
		service.logger.Warn(ctx, "failed to save user %s: %s", user.GetID(), err)
		return err
	}
	service.logger.Info(ctx, "user %s saved", user.GetID())

	service.logger.Info(ctx, "successfully registered user %s", user.GetID())
	return nil
}

// Deprecated: use read model
func (service *HumanResourcesService) GetUsers(ctx context.Context) ([]user.UserDto, error) {
	var err error

	users, err := service.userStore.GetMany(ctx)

	usersDto := make([]user.UserDto, 0)
	for _, u := range users {
		userDto := user.DtoFrom(u)
		usersDto = append(usersDto, userDto)
	}

	return usersDto, err
}

func (service *HumanResourcesService) UserJoinGroup(ctx context.Context, userId user.UserID, groupId group.GroupID) error {
	var err error
	service.logger.With(ctx, "userId", userId.String()).With(ctx, "groupId", groupId.String())
	service.logger.Info(ctx, "starting UserJoinGroup for user %s on group %s", userId, groupId)

	groupToJoin, err := service.groupStore.Get(ctx, groupId)
	if err != nil {
		service.logger.Info(ctx, "failed to get group %s: %s", groupId, err)
		return err
	}

	userToJoin, err := service.userStore.Get(ctx, userId)
	if err != nil {
		service.logger.Warn(ctx, "failed to get user %s: %s", userId, err)
		return err
	}

	err = groupToJoin.AddMember(userToJoin.GetID())
	if err != nil {
		service.logger.Info(ctx, "failed to add user %s to group %s: %s", userToJoin.GetID(), groupToJoin.GetID(), err)
		return err
	}

	err = service.groupStore.Save(ctx, groupToJoin)
	if err != nil {
		service.logger.Warn(ctx, "failed to save group %s: %s", groupToJoin.GetID(), err)
		return err
	}

	service.logger.Info(ctx, "successfully joined user %s to group %s", userToJoin.GetID(), groupToJoin.GetID())
	return nil
}
