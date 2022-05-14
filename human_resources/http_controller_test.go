//go:build intg

package human_resources

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"go-microservices/common"
	"go-microservices/human_resources/group"
	"go-microservices/human_resources/user"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestHttpRegister(t *testing.T) {
	assert := assert.New(t)
	userStore := user.NewFakeUserStore()
	groupStore := group.NewFakeGroupStore()
	privateKey, _ := rsa.GenerateKey(rand.Reader, 1024)
	publicKey := &privateKey.PublicKey
	service := NewHumanResourcesService(&userStore, &groupStore)
	controller := HumanResourcesHttpController{humanResourcesService: service, publicKey: *publicKey, rootPath: ""}

	t.Run("it should send a 401 if there is no Token", func(t *testing.T) {
		req := common.NewRequestBuilder("POST", "/register").WithPayload(UserRegisterDto{
			PreferredName: "Shaylyn Ognjan",
		}).Build()

		rr := httptest.NewRecorder()
		controller.ServeHTTP(rr, &req)

		assert.Equal(http.StatusUnauthorized, rr.Code)
		var payload common.ErrorDto
		json.NewDecoder(rr.Body).Decode(&payload)
		assert.Equal("request has no authorization header", payload.Error)
	})

	t.Run("it should send a 401 if the Token is wrong", func(t *testing.T) {
		req := common.NewRequestBuilder("POST", "/register").WithPayload(UserRegisterDto{
			PreferredName: "John Doe",
		}).Build()

		req.Header.Add("Authorization", "Bearer dummy")

		rr := httptest.NewRecorder()
		controller.ServeHTTP(rr, &req)

		assert.Equal(http.StatusUnauthorized, rr.Code)
		var payload common.ErrorDto
		json.NewDecoder(rr.Body).Decode(&payload)
		assert.Equal("error decoding token: illegal base64 data at input byte 4", payload.Error)
	})

	t.Run("it should send a 201 and register the user if the Token is valid", func(t *testing.T) {
		token := common.Token{
			CreatedAt:  time.Now(),
			Identifier: "phyliss@otto.com",
		}
		req := common.NewRequestBuilder("POST", "/register").WithPayload(UserRegisterDto{
			PreferredName: "Phyliss Ottó",
		}).WithToken(token, *privateKey).Build()

		rr := httptest.NewRecorder()
		controller.ServeHTTP(rr, &req)

		assert.Equal(http.StatusCreated, rr.Code)
		var payload common.OperationDto
		json.NewDecoder(rr.Body).Decode(&payload)
		assert.True(payload.Success)

		_, err := userStore.GetByEmail(ctx, "phyliss@otto.com")
		assert.NoError(err)
	})
}

func TestHttpGetUsers(t *testing.T) {
	assert := assert.New(t)
	userStore := user.NewFakeUserStore()
	privateKey, _ := rsa.GenerateKey(rand.Reader, 1024)
	publicKey := &privateKey.PublicKey
	service := NewHumanResourcesService(&userStore, nil)
	controller := HumanResourcesHttpController{humanResourcesService: service, publicKey: *publicKey, rootPath: ""}

	t.Run("it should send a 200 with the users", func(t *testing.T) {
		service.Register(ctx, "john@doe.com", "John Doe")
		savedUser, _ := userStore.GetByEmail(ctx, "john@doe.com")
		req := common.NewRequestBuilder("GET", "/").Build()

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

func TestHttpJoin(t *testing.T) {
	assert := assert.New(t)
	userStore := user.NewFakeUserStore()
	groupStore := group.NewFakeGroupStore()
	privateKey, _ := rsa.GenerateKey(rand.Reader, 1024)
	publicKey := &privateKey.PublicKey
	service := NewHumanResourcesService(&userStore, &groupStore)
	controller := HumanResourcesHttpController{humanResourcesService: service, publicKey: *publicKey, rootPath: ""}

	t.Run("it should send a 400 if the group_id is not a uuid", func(t *testing.T) {
		req := common.NewRequestBuilder("POST", "/join_group").WithPayload(UserJoinGroupDto{GroupID: "dummy", UserID: uuid.NewString()}).Build()

		rr := httptest.NewRecorder()
		controller.ServeHTTP(rr, &req)

		assert.Equal(http.StatusBadRequest, rr.Code)
		var payload common.ErrorDto
		json.NewDecoder(rr.Body).Decode(&payload)
		assert.Equal("invalid group id", payload.Error)
	})

	t.Run("it should send a 400 if the user_id is not a uuid", func(t *testing.T) {
		req := common.NewRequestBuilder("POST", "/join_group").WithPayload(UserJoinGroupDto{GroupID: uuid.NewString(), UserID: "dummy"}).Build()

		rr := httptest.NewRecorder()
		controller.ServeHTTP(rr, &req)

		assert.Equal(http.StatusBadRequest, rr.Code)
		var payload common.ErrorDto
		json.NewDecoder(rr.Body).Decode(&payload)
		assert.Equal("invalid user id", payload.Error)
	})

	t.Run("it should send a 201 and make the user join", func(t *testing.T) {
		userToJoin := user.New(user.NewUserPayload{Email: "john@doe.com", PreferredName: "John Doe"})
		groupToJoin := group.New("Group 1", "")
		userStore.Save(ctx, userToJoin)
		groupStore.Save(ctx, groupToJoin)
		req := common.NewRequestBuilder("POST", "/join_group").WithPayload(UserJoinGroupDto{GroupID: groupToJoin.GetID().String(), UserID: userToJoin.GetID().String()}).Build()

		rr := httptest.NewRecorder()
		controller.ServeHTTP(rr, &req)

		assert.Equal(http.StatusCreated, rr.Code)
		var payload common.OperationDto
		json.NewDecoder(rr.Body).Decode(&payload)
		assert.True(payload.Success)

		groupUpdated, _ := groupStore.Get(ctx, groupToJoin.GetID())
		assert.True(groupUpdated.IsMember(userToJoin.GetID()))
	})
}
