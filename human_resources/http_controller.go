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

	"github.com/go-chi/chi/v5"
)

type HumanResourcesHttpController struct {
	*chi.Mux
	humanResourcesService *HumanResourcesService
	publicKey             rsa.PublicKey
}

func NewHumanResourcesHttpController(humanResourcesService *HumanResourcesService, publicKey rsa.PublicKey) HumanResourcesHttpController {
	controller := HumanResourcesHttpController{
		Mux:                   chi.NewRouter(),
		humanResourcesService: humanResourcesService,
		publicKey:             publicKey,
	}

	controller.Get("/users", controller.GetUsers)
	controller.Post("/users/register", controller.Register)
	controller.Post("/groups/{groupId}/join", controller.JoinGroup)

	return controller
}

type UserRegisterDto struct {
	PreferredName string `json:"preferred_name"`
}

func (controller *HumanResourcesHttpController) Register(w http.ResponseWriter, req *http.Request) {
	var err error

	token, err := common.ExtractToken(*req, controller.publicKey)
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
	users, err := controller.humanResourcesService.userStore.GetMany(req.Context())
	if err != nil {
		common.WriteError(w, http.StatusBadRequest, err)
	}

	common.WriteResponse(w, http.StatusOK, user.UserListDto{Items: users, Total: len(users)})
}

type UserJoinGroupDto struct {
	UserID string `json:"user_id"`
}

func validateUserJoinGroupDto(groupIdParam string, input io.Reader) (user.UserID, group.GroupID, []error) {
	var errs []error
	payload := new(UserJoinGroupDto)
	err := json.NewDecoder(input).Decode(payload)
	if err != nil {
		errs = append(errs, fmt.Errorf("failed to parse input: %s", err.Error()))
		return user.UserID{}, group.GroupID{}, errs
	}

	id, err := common.ParseID(groupIdParam)
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
	userId, groupId, errs := validateUserJoinGroupDto(chi.URLParam(req, "groupId"), req.Body)
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
