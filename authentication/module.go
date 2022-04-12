package authentication

import (
	"crypto/rsa"
	"go-microservices/authentication/account"
)

func Bootstrap(privateKey rsa.PrivateKey) *AuthenticationService {
	accountStore := account.NewSqlAccountStore()
	return &AuthenticationService{
		accountStore: &accountStore,
		privateKey:   privateKey,
	}
}
