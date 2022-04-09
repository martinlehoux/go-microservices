package authentication

import (
	"bytes"
	"encoding/json"
	"go-microservices/common"
	"time"
)

type Token struct {
	CreatedAt  time.Time `json:"created_at"`
	Identifier string    `json:"identifier"`
}

func (token Token) Bytes() ([]byte, error) {
	encodedBytes := new(bytes.Buffer)
	err := json.NewEncoder(encodedBytes).Encode(token)
	if err != nil {
		return nil, err
	}
	return encodedBytes.Bytes(), nil
}

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

func (account *Account) CreateToken() Token {
	return Token{
		CreatedAt:  time.Now(),
		Identifier: account.Identifier,
	}
}

func (account *Account) ValidatePassword(hashedPassword []byte) bool {
	return bytes.Equal(account.HashedPassword, hashedPassword)
}
