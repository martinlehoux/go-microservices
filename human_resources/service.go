package human_resources

import (
	"errors"
	"go-microservices/human_resources/user"
	"log"
)

type HumanResourcesService struct {
	userStore user.UserStore
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

	user := user.NewUser(user.NewUserPayload{Email: email, PreferredName: preferredName})

	err = service.userStore.Save(user)
	if err != nil {
		log.Printf("failed to save user %s: %s", user.GetID(), err)
		return err
	}
	log.Printf("user %s saved", user.GetID())

	log.Printf("successfully registered user %s", user.GetID())
	return nil
}

func (service *HumanResourcesService) GetUsers() ([]user.User, error) {
	var err error

	users, err := service.userStore.GetMany()

	return users, err
}
