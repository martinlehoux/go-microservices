package authentication

import (
	"crypto/rsa"
	"go-microservices/authentication/account"
)

func Bootstrap(privateKey rsa.PrivateKey) *AuthenticationHttpController {
	accountStore := account.NewSqlAccountStore()
	authenticationService := NewAuthenticationService(&accountStore, privateKey)
	authenticationController := NewAuthenticationHttpController(&authenticationService)
	return &authenticationController
}
