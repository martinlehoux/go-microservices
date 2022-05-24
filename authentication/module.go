package authentication

import (
	"crypto/rsa"
	"go-microservices/authentication/account"
	"go-microservices/common"
)

func Bootstrap(logger common.Logger, privateKey rsa.PrivateKey) *AuthenticationHttpController {
	accountStore := account.NewSqlAccountStore()
	authenticationService := NewAuthenticationService(&accountStore, logger, privateKey)
	authenticationController := NewAuthenticationHttpController(&authenticationService)
	return &authenticationController
}
