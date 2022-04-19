//go:build intg

package human_resources

import (
	"crypto/rand"
	"crypto/rsa"
	"go-microservices/common"
	"go-microservices/human_resources/user"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestHttpRegister(t *testing.T) {
	assert := assert.New(t)
	userStore := user.NewFakeUserStore()
	privateKey, _ := rsa.GenerateKey(rand.Reader, 1024)
	publicKey := &privateKey.PublicKey
	service := HumanResourcesService{userStore: &userStore}
	controller := HumanResourcesHttpController{humanResourcesService: &service, publicKey: *publicKey, rootPath: ""}

	t.Run("it should send a 401 if there is no Token", func(t *testing.T) {
		req := common.PrepareRequest("POST", "/register", UserRegisterForm{
			PreferredName: "Shaylyn Ognjan",
		})

		rr := httptest.NewRecorder()
		controller.ServeHTTP(rr, &req)

		assert.Equal(http.StatusUnauthorized, rr.Code)
		body, _ := ioutil.ReadAll(rr.Body)
		assert.JSONEq(`{"error":"Unauthorized: request has no authorization header"}`, string(body))
	})

	t.Run("it should send a 401 if the Token is wrong", func(t *testing.T) {
		req := common.PrepareRequest("POST", "/register", UserRegisterForm{
			PreferredName: "John Doe",
		})

		req.Header.Add("Authorization", "Bearer dummy")

		rr := httptest.NewRecorder()
		controller.ServeHTTP(rr, &req)

		assert.Equal(http.StatusUnauthorized, rr.Code)
		body, _ := ioutil.ReadAll(rr.Body)
		assert.JSONEq(`{"error":"Unauthorized: invalid character 'd' looking for beginning of value"}`, string(body))
	})

	t.Run("it should send a 201 and register the user if the Token is valid", func(t *testing.T) {
		req := common.PrepareRequest("POST", "/register", UserRegisterForm{
			PreferredName: "Phyliss Ottó",
		})

		token := common.Token{
			CreatedAt:  time.Now(),
			Identifier: "phyliss@otto.com",
		}
		rawToken, _ := common.SignToken(token, *privateKey)
		req.Header.Add("Authorization", "Bearer "+string(rawToken))

		rr := httptest.NewRecorder()
		controller.ServeHTTP(rr, &req)

		assert.Equal(http.StatusCreated, rr.Code)
		body, _ := ioutil.ReadAll(rr.Body)
		assert.JSONEq(`{"success":true}`, string(body))

		_, err := userStore.GetByEmail("phyliss@otto.com")
		assert.NoError(err)
	})
}
