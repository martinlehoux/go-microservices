package authentication

import (
	"context"
	"crypto/rsa"
	"errors"
	"go-microservices/authentication/account"
	"go-microservices/common"

	"golang.org/x/crypto/argon2"
)

type AuthenticationService struct {
	accountStore account.AccountStore
	logger       common.Logger
	privateKey   rsa.PrivateKey
}

var (
	ErrWrongPassword       = errors.New("wrong password")
	ErrTokenSigningFailure = errors.New("token signing failure")
	ErrIdentifierUsed      = errors.New("identifier already used")
)

func NewAuthenticationService(accountStore account.AccountStore, logger common.Logger, privateKey rsa.PrivateKey) AuthenticationService {
	return AuthenticationService{
		accountStore: accountStore,
		logger:       logger,
		privateKey:   privateKey,
	}
}

func (service *AuthenticationService) Authenticate(ctx context.Context, identifier string, password string) ([]byte, error) {
	service.logger.With(ctx, "identifier", identifier)
	service.logger.Info(ctx, "starting authentication for identifier %s", identifier)

	account, err := service.accountStore.GetByIdentifier(ctx, identifier)
	if err != nil {
		service.logger.Info(ctx, "failed to find account for identifier %s: %s", identifier, err)
		return nil, err
	}
	service.logger.With(ctx, "accountId", account.GetID().String())

	hashedPassword := service.hashPassword(password)
	if !account.ValidatePassword(hashedPassword) {
		service.logger.Info(ctx, "failed to authenticate account for identifier %s: password mismatch", identifier)
		return nil, ErrWrongPassword
	}

	token := account.CreateToken()

	signedToken, err := common.SignToken(token, service.privateKey)
	if err != nil {
		service.logger.Warn(ctx, "failed to sign token for identifier %s: %s", identifier, err)
		return nil, ErrTokenSigningFailure
	}

	return signedToken, nil
}

func (service *AuthenticationService) Register(ctx context.Context, identifier string, password string) error {
	var err error
	service.logger.Info(ctx, "starting registration for identifier %s", identifier)

	err = service.ensureIdentifierNotUsed(ctx, identifier)
	if err != nil {
		service.logger.Info(ctx, "failed to ensure identifier %s is not used: %s", identifier, err)
		return err
	}

	hashedPassword := service.hashPassword(password)

	account := account.NewAccount(identifier, hashedPassword)
	service.logger.With(ctx, "accountId", account.GetID().String())

	err = service.accountStore.Save(ctx, account)
	if err != nil {
		service.logger.Warn(ctx, "failed to save account %s: %s", account.GetID(), err)
		return err
	}

	service.logger.Info(ctx, "successfully registered account %s", account.GetID())
	return nil
}

func (service *AuthenticationService) hashPassword(password string) []byte {
	return argon2.IDKey([]byte(password), []byte("salt"), 1, 64*1024, 4, 32)
}

func (service *AuthenticationService) ensureIdentifierNotUsed(ctx context.Context, identifier string) error {
	var err error

	_, err = service.accountStore.GetByIdentifier(ctx, identifier)

	if err == nil {
		return ErrIdentifierUsed
	}

	return nil
}
