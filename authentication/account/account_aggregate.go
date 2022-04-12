package account

import (
	"bytes"
	"go-microservices/common"
	"time"
)

type AccountID common.ID

type Account struct {
	Id             AccountID
	Identifier     string
	HashedPassword []byte
}

func NewAccount(identifier string, hashedPassword []byte) Account {
	return Account{
		Id:             AccountID(common.CreateID()),
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
