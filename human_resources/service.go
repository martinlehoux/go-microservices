package human_resources

import (
	"errors"
	"go-microservices/human_resources/group"
	"go-microservices/human_resources/user"
	"log"
)

type HumanResourcesService struct {
	userStore  user.UserStore
	groupStore group.GroupStore
}

func NewHumanResourcesService(userStore user.UserStore, groupStore group.GroupStore) *HumanResourcesService {
	return &HumanResourcesService{userStore: userStore, groupStore: groupStore}
}

func (service *HumanResourcesService) Register(email string, preferredName string) error {
	var err error
	log.Printf("starting user registration for email %s", email)

	emailUsed, err := service.userStore.EmailExists(email)
	if err != nil {
		log.Printf("failed to check if email %s exists: %s", email, err)
		return err
	}
	if emailUsed {
		err = errors.New("email already used")
		log.Printf("failed to register user for email %s: %s", email, err)
		return err
	}

	user := user.New(user.NewUserPayload{Email: email, PreferredName: preferredName})

	err = service.userStore.Save(user)
	if err != nil {
		log.Printf("failed to save user %s: %s", user.GetID(), err)
		return err
	}
	log.Printf("user %s saved", user.GetID())

	log.Printf("successfully registered user %s", user.GetID())
	return nil
}

func (service *HumanResourcesService) GetUsers() ([]user.UserDto, error) {
	var err error

	users, err := service.userStore.GetMany()

	usersDto := make([]user.UserDto, 0)
	for _, u := range users {
		userDto := user.DtoFrom(u)
		usersDto = append(usersDto, userDto)
	}

	return usersDto, err
}

func (service *HumanResourcesService) UserJoinGroup(userId user.UserID, groupId group.GroupID) error {
	var err error

	groupToJoin, err := service.groupStore.Get(groupId)
	if err != nil {
		log.Printf("failed to get group %s: %s", groupId, err)
		return err
	}

	userToJoin, err := service.userStore.Get(userId)
	if err != nil {
		log.Printf("failed to get user %s: %s", userId.String(), err)
		return err
	}

	err = groupToJoin.AddMember(userToJoin.GetID())
	if err != nil {
		log.Printf("failed to add user %s to group %s: %s", userToJoin.GetID(), groupToJoin.GetID(), err)
		return err
	}

	err = service.groupStore.Save(groupToJoin)
	if err != nil {
		log.Printf("failed to save group %s: %s", groupToJoin.GetID(), err)
		return err
	}

	return nil
}
