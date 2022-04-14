package common

import (
	"crypto/rand"
	"crypto/rsa"
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
		token := Token{CreatedAt: time.Now(), Identifier: "test"}
		blob, _ := SignToken(token, *randomKey)

		token, err := ParseToken(blob, publicKey)

		assert.ErrorContains(err, "verification error", "the error should be returned")
		assert.Equal(token, Token{}, "the token should be empty")
	})

	t.Run("it should error for an invalid signature", func(t *testing.T) {
		signature := make([]byte, 32)
		rand.Read(signature)
		token := Token{CreatedAt: time.Now(), Identifier: "test", Signature: signature}
		blob, _ := token.Bytes()

		token, err := ParseToken(blob, publicKey)

		assert.ErrorContains(err, "verification error", "the error should be returned")
		assert.Equal(token, Token{}, "the token should be empty")
	})

	t.Run("it should parse correctly a signed token", func(t *testing.T) {
		token := Token{CreatedAt: time.Now(), Identifier: "test"}
		blob, _ := SignToken(token, *privateKey)

		parsedToken, err := ParseToken(blob, publicKey)

		assert.NoError(err, "the parsing should be successful")
		assert.True(token.CreatedAt.Equal(parsedToken.CreatedAt), "the creation date should match")
		assert.Equal(token.Identifier, parsedToken.Identifier, "the identifier should match")
	})

	t.Run("it should error for an expired Token", func(t *testing.T) {
		token := Token{CreatedAt: time.Now().Add(-2 * time.Hour), Identifier: "test"}
		blob, _ := SignToken(token, *privateKey)

		token, err := ParseToken(blob, publicKey)

		assert.ErrorContains(err, "token is expired", "the error should be returned")
		assert.Equal(token, Token{}, "the token should be empty")
	})
}
