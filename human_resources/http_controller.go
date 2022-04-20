package human_resources

import (
	"crypto/rsa"
	"encoding/json"
	"go-microservices/common"
	"go-microservices/human_resources/user"
	"net/http"
	"strings"
)

type HumanResourcesHttpController struct {
	humanResourcesService *HumanResourcesService
	publicKey             rsa.PublicKey
	rootPath              string
}

func (controller *HumanResourcesHttpController) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	path := strings.TrimPrefix(req.URL.Path, controller.rootPath)
	if path == "/register" && req.Method == "POST" {
		controller.Register(w, req)
	}
	if path == "/" && req.Method == "GET" {
		controller.GetUsers(w, req)
	}
}

type UserRegisterDto struct {
	PreferredName string `json:"preferred_name"`
}

func (controller *HumanResourcesHttpController) Register(w http.ResponseWriter, req *http.Request) {
	var err error

	token, err := common.ExtractToken(req.Header.Get("Authorization"), controller.publicKey)
	if err != nil {
		common.WriteError(w, http.StatusUnauthorized, err)
		return
	}

	payload := new(UserRegisterDto)
	err = json.NewDecoder(req.Body).Decode(payload)
	if err != nil {
		common.WriteError(w, http.StatusBadRequest, err)
		return
	}

	err = controller.humanResourcesService.Register(token.Identifier, payload.PreferredName)
	if err != nil {
		common.WriteError(w, http.StatusBadRequest, err)
		return
	}

	common.WriteResponse(w, http.StatusCreated, common.OperationDto{Success: true})
}

func (controller *HumanResourcesHttpController) GetUsers(w http.ResponseWriter, req *http.Request) {
	users, err := controller.humanResourcesService.GetUsers()
	if err != nil {
		common.WriteError(w, http.StatusBadRequest, err)
	}

	common.WriteResponse(w, http.StatusOK, user.UserListDto{Items: users, Total: len(users)})
}
