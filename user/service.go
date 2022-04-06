package user

import (
	"errors"
	"go-microservices/authentication"
	"log"
)

type UserService struct {
	userStore             UserStore
	authenticationService *authentication.AuthenticationService
}

func Bootstrap(authenticationService *authentication.AuthenticationService) *UserService {
	store := bootstrapFakeUserStore()
	return &UserService{
		userStore:             &store,
		authenticationService: authenticationService,
	}
}

func (service *UserService) Register(email string) error {
	var err error
	log.Printf("starting user registration for email %s", email)

	emailUsed, err := service.userStore.emailExists(email)
	if err != nil {
		log.Printf("failed to check if email %s exists: %s", email, err)
		return err
	}
	if emailUsed {
		err = errors.New("email already used")
		log.Printf("failed to register user for email %s: %s", email, err)
		return err
	}

	user := NewUser(NewUserPayload{Email: email, PreferredName: ""})

	err = service.userStore.save(user)
	if err != nil {
		log.Printf("failed to save user %s: %s", user.Id, err)
		return err
	}
	log.Printf("user %s saved", user.Id)

	log.Printf("successfully registered user %s", user.Id)
	return nil
}

func (service *UserService) GetUsers() ([]User, error) {
	var err error

	users, err := service.userStore.getMany()

	return users, err
}
