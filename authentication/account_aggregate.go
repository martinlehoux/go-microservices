package authentication

import (
	"bytes"
	"go-microservices/common"
	"time"
)

type Account struct {
	Id             common.ID
	Identifier     string
	HashedPassword []byte
}

func NewAccount(identifier string, hashedPassword []byte) Account {
	return Account{
		Id:             common.CreateID(),
		Identifier:     identifier,
		HashedPassword: hashedPassword,
	}
}

func (account *Account) CreateToken() common.Token {
	return common.Token{
		CreatedAt:  time.Now(),
		Identifier: account.Identifier,
	}
}

func (account *Account) ValidatePassword(hashedPassword []byte) bool {
	return bytes.Equal(account.HashedPassword, hashedPassword)
}
