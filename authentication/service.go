package authentication

import (
	"context"
	"crypto/rsa"
	"errors"
	"go-microservices/authentication/account"
	"go-microservices/common"
	"log"

	"golang.org/x/crypto/argon2"
)

type AuthenticationService struct {
	accountStore account.AccountStore
	privateKey   rsa.PrivateKey
}

var (
	ErrWrongPassword       = errors.New("wrong password")
	ErrTokenSigningFailure = errors.New("token signing failure")
	ErrIdentifierUsed      = errors.New("identifier already used")
)

func NewAuthenticationService(accountStore account.AccountStore, privateKey rsa.PrivateKey) AuthenticationService {
	return AuthenticationService{
		accountStore: accountStore,
		privateKey:   privateKey,
	}
}

func (service *AuthenticationService) Authenticate(ctx context.Context, identifier string, password string) ([]byte, error) {
	log.Printf("starting authentication for identifier %s", identifier)

	account, err := service.accountStore.GetByIdentifier(ctx, identifier)
	if err != nil {
		log.Printf("failed to find account for identifier %s: %s", identifier, err)
		return nil, err
	}

	hashedPassword := service.hashPassword(password)
	if !account.ValidatePassword(hashedPassword) {
		log.Printf("failed to authenticate account for identifier %s: password mismatch", identifier)
		return nil, ErrWrongPassword
	}

	token := account.CreateToken()

	signedToken, err := common.SignToken(token, service.privateKey)
	if err != nil {
		log.Printf("failed to sign token for identifier %s: %s", identifier, err)
		return nil, ErrTokenSigningFailure
	}

	return signedToken, nil
}

func (service *AuthenticationService) Register(ctx context.Context, identifier string, password string) error {
	var err error
	log.Printf("starting registration for identifier %s", identifier)

	err = service.ensureIdentifierNotUsed(ctx, identifier)
	if err != nil {
		log.Printf("failed to ensure identifier %s is not used: %s", identifier, err)
		return err
	}

	hashedPassword := service.hashPassword(password)

	account := account.NewAccount(identifier, hashedPassword)

	err = service.accountStore.Save(ctx, account)
	if err != nil {
		log.Printf("failed to save account %s: %s", account.GetID(), err)
		return err
	}

	log.Printf("successfully registered account %s", account.GetID())
	return nil
}

func (service *AuthenticationService) hashPassword(password string) []byte {
	return argon2.IDKey([]byte(password), []byte("salt"), 1, 64*1024, 4, 32)
}

func (service *AuthenticationService) ensureIdentifierNotUsed(ctx context.Context, identifier string) error {
	var err error
	log.Printf("starting check for unused identifier %s", identifier)

	_, err = service.accountStore.GetByIdentifier(ctx, identifier)

	if err == nil {
		return ErrIdentifierUsed
	}

	log.Printf("successfully ensured identifier %s is not used", identifier)
	return nil
}
