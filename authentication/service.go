package authentication

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha512"
	"errors"
	"go-microservices/common"
	"log"

	"golang.org/x/crypto/argon2"
)

type AuthenticationService struct {
	accountStore AccountStore
	privateKey   rsa.PrivateKey
}

func Bootstrap(privateKey rsa.PrivateKey) *AuthenticationService {
	accountStore := bootstrapFakeAccountStore()
	return &AuthenticationService{
		accountStore: &accountStore,
		privateKey:   privateKey,
	}
}

func (service *AuthenticationService) Authenticate(identifier string, password string) (token common.Token, signature []byte, err error) {
	log.Printf("starting authentication for identifier %s", identifier)

	account, err := service.accountStore.loadForIdentifier(identifier)
	if err != nil {
		log.Printf("failed to find account for identifier %s: %s", identifier, err)
		return common.Token{}, nil, err
	}

	hashedPassword := service.hashPassword(password)
	if !account.ValidatePassword(hashedPassword) {
		log.Printf("failed to authenticate account for identifier %s: password mismatch", identifier)
		return common.Token{}, nil, errors.New("password mismatch")
	}

	token = account.CreateToken()

	signedToken, err := service.signToken(token)
	if err != nil {
		log.Printf("failed to sign token for identifier %s: %s", identifier, err)
		return common.Token{}, nil, errors.New("failed to sign token")
	}

	return token, signedToken, nil
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

func (service *AuthenticationService) signToken(token common.Token) ([]byte, error) {
	msg, err := token.Bytes()
	if err != nil {
		return nil, err
	}
	digest := sha512.Sum512(msg)
	return rsa.SignPKCS1v15(rand.Reader, &service.privateKey, crypto.SHA512, digest[:])
}
