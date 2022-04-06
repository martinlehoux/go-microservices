package authentication

import (
	"bytes"
	"go-microservices/common"
	"time"
)

type Token struct {
	Id        common.ID
	CreatedAt time.Time
}

type Account struct {
	Id             common.ID
	Identifier     string
	HashedPassword []byte
	Tokens         []Token
}

func NewAccount(identifier string, hashedPassword []byte) Account {
	return Account{
		Id:             common.CreateID(),
		Identifier:     identifier,
		HashedPassword: hashedPassword,
	}
}

func (account *Account) CreateToken() {
	token := Token{
		Id:        common.CreateID(),
		CreatedAt: time.Now(),
	}
	account.Tokens = append(account.Tokens, token)
}

func (account *Account) ValidatePassword(hashedPassword []byte) bool {
	return bytes.Equal(account.HashedPassword, hashedPassword)
}
