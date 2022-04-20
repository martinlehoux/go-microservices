//go:build intg

package human_resources

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"go-microservices/common"
	"go-microservices/human_resources/user"
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
		req := common.PrepareRequest("POST", "/register", UserRegisterDto{
			PreferredName: "Shaylyn Ognjan",
		})

		rr := httptest.NewRecorder()
		controller.ServeHTTP(rr, &req)

		assert.Equal(http.StatusUnauthorized, rr.Code)
		var payload common.ErrorDto
		json.NewDecoder(rr.Body).Decode(&payload)
		assert.Equal("request has no authorization header", payload.Error)
	})

	t.Run("it should send a 401 if the Token is wrong", func(t *testing.T) {
		req := common.PrepareRequest("POST", "/register", UserRegisterDto{
			PreferredName: "John Doe",
		})

		req.Header.Add("Authorization", "Bearer dummy")

		rr := httptest.NewRecorder()
		controller.ServeHTTP(rr, &req)

		assert.Equal(http.StatusUnauthorized, rr.Code)
		var payload common.ErrorDto
		json.NewDecoder(rr.Body).Decode(&payload)
		assert.Equal("error decoding token: illegal base64 data at input byte 4", payload.Error)
	})

	t.Run("it should send a 201 and register the user if the Token is valid", func(t *testing.T) {
		req := common.PrepareRequest("POST", "/register", UserRegisterDto{
			PreferredName: "Phyliss Ottó",
		})

		token := common.Token{
			CreatedAt:  time.Now(),
			Identifier: "phyliss@otto.com",
		}
		rawToken, _ := common.SignToken(token, *privateKey)
		req.Header.Add("Authorization", "Bearer "+base64.StdEncoding.EncodeToString(rawToken))

		rr := httptest.NewRecorder()
		controller.ServeHTTP(rr, &req)

		assert.Equal(http.StatusCreated, rr.Code)
		var payload common.OperationDto
		json.NewDecoder(rr.Body).Decode(&payload)
		assert.True(payload.Success)

		_, err := userStore.GetByEmail("phyliss@otto.com")
		assert.NoError(err)
	})
}

func TestHttpGetUsers(t *testing.T) {
	assert := assert.New(t)
	userStore := user.NewFakeUserStore()
	privateKey, _ := rsa.GenerateKey(rand.Reader, 1024)
	publicKey := &privateKey.PublicKey
	service := HumanResourcesService{userStore: &userStore}
	controller := HumanResourcesHttpController{humanResourcesService: &service, publicKey: *publicKey, rootPath: ""}

	t.Run("it should send a 200 with the users", func(t *testing.T) {
		service.Register("john@doe.com", "John Doe")
		savedUser, _ := userStore.GetByEmail("john@doe.com")
		req := common.PrepareRequest("GET", "/", nil)

		rr := httptest.NewRecorder()
		controller.ServeHTTP(rr, &req)

		assert.Equal(http.StatusOK, rr.Code)
		var payload user.UserListDto
		json.NewDecoder(rr.Body).Decode(&payload)
		assert.Equal(1, payload.Total)
		assert.Equal(user.UserListDto{
			Items: []user.UserDto{
				{ID: savedUser.GetID().String(), PreferredName: "John Doe", Email: "john@doe.com"},
			},
			Total: 1,
		}, payload)
	})
}
