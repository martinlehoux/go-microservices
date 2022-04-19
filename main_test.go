//go:build e2e

package main

import (
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPing(t *testing.T) {
	assert := assert.New(t)

	t.Run("it should respond with 'pong'", func(t *testing.T) {
		res, err := http.DefaultClient.Get("http://localhost:3000/ping")

		assert.NoError(err)
		assert.Equal(http.StatusOK, res.StatusCode)
		body, _ := ioutil.ReadAll(res.Body)
		assert.JSONEq(`{"message":"pong"}`, string(body))
	})
}
