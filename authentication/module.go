package authentication

import (
	"crypto/rsa"
)

func Bootstrap(privateKey rsa.PrivateKey) *AuthenticationService {
	accountStore := NewSqlAccountStore()
	return &AuthenticationService{
		accountStore: &accountStore,
		privateKey:   privateKey,
	}
}
