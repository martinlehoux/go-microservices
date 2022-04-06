package authentication

import (
	"errors"
	"log"

	"golang.org/x/crypto/argon2"
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

func (service *AuthenticationService) Authenticate(identifier string, password string) (Token, error) {
	var err error
	log.Printf("starting authentication for identifier %s", identifier)

	account, err := service.accountStore.loadForIdentifier(identifier)
	if err != nil {
		log.Printf("failed to find account for identifier %s: %s", identifier, err)
		return Token{}, err
	}

	hashedPassword := service.hashPassword(password)
	if !account.ValidatePassword(hashedPassword) {
		log.Printf("failed to authenticate account for identifier %s: password mismatch", identifier)
		return Token{}, errors.New("password mismatch")
	}

	account.CreateToken()

	err = service.accountStore.save(account)
	if err != nil {
		log.Printf("failed to save account %s: %s", account.Id, err)
		return Token{}, err
	}

	log.Printf("successfully authenticated identifier %s", identifier)
	return account.Tokens[len(account.Tokens)-1], nil
}

func (service *AuthenticationService) Register(identifier string, password string) error {
	var err error
	log.Printf("starting registration for identifier %s", identifier)

	err = service.ensureIdentifierNotUsed(identifier)
	if err != nil {
		log.Printf("failed to ensure identifier %s is not used: %s", identifier, err)
		return err
	}

	hashedPassword := service.hashPassword(password)

	account := NewAccount(identifier, hashedPassword)

	err = service.accountStore.save(account)
	if err != nil {
		log.Printf("failed to save account %s: %s", account.Id, err)
		return err
	}

	log.Printf("successfully registered account %s", account.Id)
	return nil
}

func (service *AuthenticationService) hashPassword(password string) []byte {
	return argon2.IDKey([]byte(password), []byte("salt"), 1, 64*1024, 4, 32)
}

func (service *AuthenticationService) ensureIdentifierNotUsed(identifier string) error {
	var err error
	log.Printf("starting check for unused identifier %s", identifier)

	_, err = service.accountStore.loadForIdentifier(identifier)

	if err == nil {
		return errors.New("identifier already used")
	}

	log.Printf("successfully ensured identifier %s is not used", identifier)
	return nil
}
