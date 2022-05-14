package human_resources

import (
	"crypto/rsa"
	"encoding/json"
	"errors"
	"go-microservices/common"
	"go-microservices/human_resources/group"
	"go-microservices/human_resources/user"
	"net/http"
	"strings"
)

type HumanResourcesHttpController struct {
	humanResourcesService *HumanResourcesService
	publicKey             rsa.PublicKey
	rootPath              string
}

func NewHumanResourcesHttpController(humanResourcesService *HumanResourcesService, publicKey rsa.PublicKey, rootPath string) HumanResourcesHttpController {
	return HumanResourcesHttpController{
		humanResourcesService: humanResourcesService,
		publicKey:             publicKey,
		rootPath:              rootPath,
	}
}

func (controller *HumanResourcesHttpController) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	path := strings.TrimPrefix(req.URL.Path, controller.rootPath)
	if path == "/register" && req.Method == "POST" {
		controller.Register(w, req)
	} else if path == "/" && req.Method == "GET" {
		controller.GetUsers(w, req)
	} else if path == "/join_group" && req.Method == "POST" {
		controller.JoinGroup(w, req)
	} else {
		common.WriteError(w, http.StatusNotFound, errors.New("not found"))
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

	err = controller.humanResourcesService.Register(req.Context(), token.Identifier, payload.PreferredName)
	if err != nil {
		common.WriteError(w, http.StatusBadRequest, err)
		return
	}

	common.WriteResponse(w, http.StatusCreated, common.OperationDto{Success: true})
}

func (controller *HumanResourcesHttpController) GetUsers(w http.ResponseWriter, req *http.Request) {
	users, err := controller.humanResourcesService.GetUsers(req.Context())
	if err != nil {
		common.WriteError(w, http.StatusBadRequest, err)
	}

	common.WriteResponse(w, http.StatusOK, user.UserListDto{Items: users, Total: len(users)})
}

type UserJoinGroupDto struct {
	GroupID string `json:"group_id"`
	UserID  string `json:"user_id"`
}

func (controller *HumanResourcesHttpController) JoinGroup(w http.ResponseWriter, req *http.Request) {
	payload := new(UserJoinGroupDto)
	err := json.NewDecoder(req.Body).Decode(payload)
	if err != nil {
		common.WriteError(w, http.StatusBadRequest, err)
		return
	}

	id, err := common.ParseID(payload.GroupID)
	if err != nil {
		common.WriteError(w, http.StatusBadRequest, errors.New("invalid group id"))
		return
	}
	groupId := group.GroupID{id}

	id, err = common.ParseID(payload.UserID)
	if err != nil {
		common.WriteError(w, http.StatusBadRequest, errors.New("invalid user id"))
		return
	}
	userId := user.UserID{id}

	err = controller.humanResourcesService.UserJoinGroup(req.Context(), userId, groupId)
	if err != nil {
		common.WriteError(w, http.StatusUnprocessableEntity, err)
		return
	}

	common.WriteResponse(w, http.StatusCreated, common.OperationDto{Success: true})
}
