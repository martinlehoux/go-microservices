package authentication

import (
	"go-microservices/common"
	"log"
)

type AuthenticationService struct {
	accountStore AccountStore
}

func Bootstrap() *AuthenticationService {
	accountStore := bootstrapFakeAccountStore()
	return &AuthenticationService{
		accountStore: &accountStore,
	}
}

func (service *AuthenticationService) CreateAccount(targetId common.ID) error {
	var err error
	log.Printf("starting account creation for target %s", targetId)

	account := NewAccount(targetId)

	err = service.accountStore.save(account)
	if err != nil {
		log.Printf("failed to save account %s: %s", account.Id, err)
		return err
	}
	log.Printf("account %s saved for target %s", targetId, account.Id)

	log.Printf("successfully created account %s", account.Id)
	return nil
}
