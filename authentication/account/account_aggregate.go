package account

import (
	"bytes"
	"go-microservices/common"
	"time"
)

type AccountID struct{ common.ID }

type Account struct {
	id             AccountID
	identifier     string
	hashedPassword []byte
}

func NewAccount(identifier string, hashedPassword []byte) Account {
	return Account{
		id:             AccountID{common.CreateID()},
		identifier:     identifier,
		hashedPassword: hashedPassword,
	}
}

func (account *Account) CreateToken() common.Token {
	return common.Token{
		CreatedAt:  time.Now(),
		Identifier: account.identifier,
	}
}

func (account *Account) ValidatePassword(hashedPassword []byte) bool {
	return bytes.Equal(account.hashedPassword, hashedPassword)
}

func (account *Account) GetID() AccountID {
	return account.id
}
