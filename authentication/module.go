package authentication

import (
	"crypto/rsa"
	"go-microservices/authentication/account"
)

func Bootstrap(rootPath string, privateKey rsa.PrivateKey) *AuthenticationHttpController {
	accountStore := account.NewSqlAccountStore()
	authenticationService := AuthenticationService{
		accountStore: &accountStore,
		privateKey:   privateKey,
	}
	authenticationController := AuthenticationHttpController{
		authenticationService: &authenticationService,
		rootPath:              rootPath,
	}
	return &authenticationController
}
