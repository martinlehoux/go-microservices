package human_resources

import (
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"go-microservices/common"
	"go-microservices/human_resources/group"
	"go-microservices/human_resources/user"
	"io"
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
		common.WriteError(w, http.StatusNotFound, common.ErrURLNotFound)
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

func validateUserJoinGroupDto(input io.Reader) (user.UserID, group.GroupID, []error) {
	var errs []error
	payload := new(UserJoinGroupDto)
	err := json.NewDecoder(input).Decode(payload)
	if err != nil {
		errs = append(errs, fmt.Errorf("failed to parse input: %s", err.Error()))
		return user.UserID{}, group.GroupID{}, errs
	}

	id, err := common.ParseID(payload.GroupID)
	if err != nil {
		errs = append(errs, fmt.Errorf("invalid group id: %s", err.Error()))
	}
	groupId := group.GroupID{id}

	id, err = common.ParseID(payload.UserID)
	if err != nil {
		errs = append(errs, fmt.Errorf("invalid user id: %s", err.Error()))
	}
	userId := user.UserID{id}

	return userId, groupId, errs
}

func (controller *HumanResourcesHttpController) JoinGroup(w http.ResponseWriter, req *http.Request) {
	userId, groupId, errs := validateUserJoinGroupDto(req.Body)
	if len(errs) > 0 {
		common.WriteErrors(w, http.StatusBadRequest, errs)
		return
	}

	err := controller.humanResourcesService.UserJoinGroup(req.Context(), userId, groupId)
	if err != nil {
		common.WriteErrors(w, http.StatusUnprocessableEntity, []error{err})
		return
	}

	common.WriteResponse(w, http.StatusCreated, common.OperationDto{Success: true})
}
