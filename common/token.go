package common

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha512"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"os"
	"time"
)

type Token struct {
	CreatedAt  time.Time `json:"created_at"`
	Identifier string    `json:"identifier"`
	Signature  []byte    `json:"signature"`
}

const TokenDuration = time.Hour

func (token Token) Bytes() ([]byte, error) {
	encodedBytes := new(bytes.Buffer)
	err := json.NewEncoder(encodedBytes).Encode(token)
	if err != nil {
		return nil, err
	}
	return encodedBytes.Bytes(), nil
}

func LoadPrivateKey(filename string) rsa.PrivateKey {
	privatePem, err := os.ReadFile(filename)
	PanicOnError(err)
	privateBlock, _ := pem.Decode(privatePem)
	print(privateBlock.Bytes)
	privateKey, err := x509.ParsePKCS1PrivateKey(privateBlock.Bytes)
	PanicOnError(err)
	return *privateKey
}

func verify(token Token, publicKey rsa.PublicKey) error {
	signature := token.Signature
	token.Signature = nil
	msg, err := token.Bytes()
	if err != nil {
		return err
	}
	digest := sha512.Sum512(msg)
	err = rsa.VerifyPKCS1v15(&publicKey, crypto.SHA512, digest[:], signature)
	return err
}

func ParseToken(blob []byte, publicKey rsa.PublicKey) (Token, error) {
	var token Token
	err := json.Unmarshal(blob, &token)
	if err != nil {
		return token, err
	}
	err = verify(token, publicKey)
	if err != nil {
		return Token{}, err
	}
	if time.Now().After(token.CreatedAt.Add(TokenDuration)) {
		return Token{}, errors.New("token is expired")
	}
	return token, nil
}

func SignToken(token Token, privateKey rsa.PrivateKey) ([]byte, error) {
	msg, err := token.Bytes()
	if err != nil {
		return nil, err
	}
	digest := sha512.Sum512(msg)
	signature, err := rsa.SignPKCS1v15(rand.Reader, &privateKey, crypto.SHA512, digest[:])
	if err != nil {
		return nil, err
	}
	token.Signature = signature
	return token.Bytes()
}
