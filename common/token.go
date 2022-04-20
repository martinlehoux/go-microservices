package common

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha512"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"os"
	"strings"
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
		return nil, fmt.Errorf("error encoding token: %s", err.Error())
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

// ExtractToken uses ParseToken on a request headers to auhtenticate a request
// TODO: After moving to het/http, extract on request instead of header string
func ExtractToken(authorization string, publicKey rsa.PublicKey) (Token, error) {
	// authorization := req.Header.Get("Authorization")
	if authorization == "" {
		return Token{}, errors.New("request has no authorization header")
	}

	bearer := strings.Split(authorization, " ")
	if len(bearer) != 2 || bearer[0] != "Bearer" {
		return Token{}, errors.New("request has invalid authorization header")
	}

	decoded, err := base64.StdEncoding.DecodeString(bearer[1])
	if err != nil {
		return Token{}, fmt.Errorf("error decoding token: %s", err.Error())
	}

	token, err := ParseToken(decoded, publicKey)
	if err != nil {
		return Token{}, fmt.Errorf("invalid token: %s", err.Error())
	}

	return token, nil
}
