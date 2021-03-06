//go:build spec

package common

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestParseToken(t *testing.T) {
	assert := assert.New(t)
	privateKey, _ := rsa.GenerateKey(rand.Reader, 1024)
	publicKey := privateKey.PublicKey

	t.Run("it should error for an invalid token", func(t *testing.T) {
		token, err := ParseToken([]byte("ABC"), publicKey)

		assert.ErrorContains(err, "invalid character", "the error should be returned")
		assert.Equal(token, Token{}, "the token should be empty")
	})

	t.Run("it should error for an invalid key", func(t *testing.T) {
		randomKey, _ := rsa.GenerateKey(rand.Reader, 1024)
		token := NewToken("test")
		blob, _ := SignToken(token, *randomKey)

		token, err := ParseToken(blob, publicKey)

		assert.ErrorContains(err, "verification error", "the error should be returned")
		assert.Equal(token, Token{}, "the token should be empty")
	})

	t.Run("it should error for an invalid signature", func(t *testing.T) {
		token := NewToken("test")
		token.Signature = make([]byte, 32)
		rand.Read(token.Signature)
		blob, _ := token.Bytes()

		token, err := ParseToken(blob, publicKey)

		assert.ErrorContains(err, "verification error", "the error should be returned")
		assert.Equal(token, Token{}, "the token should be empty")
	})

	t.Run("it should parse correctly a signed token", func(t *testing.T) {
		token := NewToken("test")
		blob, _ := SignToken(token, *privateKey)

		parsedToken, err := ParseToken(blob, publicKey)

		assert.NoError(err, "the parsing should be successful")
		assert.True(token.CreatedAt.Equal(parsedToken.CreatedAt), "the creation date should match")
		assert.Equal(token.Identifier, parsedToken.Identifier, "the identifier should match")
	})

	t.Run("it should error for an expired Token", func(t *testing.T) {
		token := NewToken("test")
		token.CreatedAt = token.CreatedAt.Add(-2 * time.Hour)
		blob, _ := SignToken(token, *privateKey)

		token, err := ParseToken(blob, publicKey)

		assert.ErrorIs(err, ErrTokenExpired)
		assert.Equal(token, Token{}, "the token should be empty")
	})
}

func TestExtractToken(t *testing.T) {
	assert := assert.New(t)
	privateKey, _ := rsa.GenerateKey(rand.Reader, 1024)
	publicKey := privateKey.PublicKey

	t.Run("it should error if there is no Authorization header", func(t *testing.T) {
		req := NewRequestBuilder("POST", "/any").Build()

		token, err := ExtractToken(req, publicKey)
		assert.ErrorIs(err, ErrMissingAuthorizationHeader)
		assert.Equal(Token{}, token)
	})

	t.Run("it should error if the header is missformed", func(t *testing.T) {
		req := NewRequestBuilder("POST", "/any").Build()
		req.Header.Add("Authorization", "Whatever")

		token, err := ExtractToken(req, publicKey)
		assert.ErrorIs(err, ErrInvalidAuthorizationHeader)
		assert.Equal(Token{}, token)
	})

	t.Run("it should error if the token is invalid", func(t *testing.T) {
		req := NewRequestBuilder("POST", "/any").Build()
		req.Header.Add("Authorization", "Bearer ABC")

		token, err := ExtractToken(req, publicKey)
		assert.ErrorContains(err, "error decoding token")
		assert.Equal(Token{}, token)
	})

	t.Run("it should return the token if it is valid", func(t *testing.T) {
		token := NewToken("test")
		blob, _ := SignToken(token, *privateKey)
		req := NewRequestBuilder("POST", "/any").Build()
		encoded := base64.StdEncoding.EncodeToString(blob)
		req.Header.Add("Authorization", "Bearer "+encoded)

		parsedToken, err := ExtractToken(req, publicKey)
		assert.NoError(err, "the parsing should be successful")
		assert.True(token.CreatedAt.Equal(parsedToken.CreatedAt), "the creation date should match")
		assert.Equal(token.Identifier, parsedToken.Identifier, "the identifier should match")
	})

	t.Run("it should error if the token is expired", func(t *testing.T) {
		t.Skip()
	})
}
