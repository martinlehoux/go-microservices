package common

import (
	"bytes"
	"crypto"
	"crypto/rsa"
	"crypto/sha512"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"os"
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

func LoadPrivateKey(filename string) rsa.PrivateKey {
	privatePem, err := os.ReadFile(filename)
	PanicOnError(err)
	privateBlock, _ := pem.Decode(privatePem)
	print(privateBlock.Bytes)
	privateKey, err := x509.ParsePKCS1PrivateKey(privateBlock.Bytes)
	PanicOnError(err)
	return *privateKey
}

func ValidateToken(publicKey rsa.PublicKey, token Token, signature []byte) error {
	msg, err := token.Bytes()
	if err != nil {
		return err
	}
	digest := sha512.Sum512(msg)
	err = rsa.VerifyPKCS1v15(&publicKey, crypto.SHA512, digest[:], signature)
	return err
}
