package authentication

import (
	"encoding/json"
	"go-microservices/common"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type AuthenticationHttpController struct {
	*chi.Mux
	authenticationService *AuthenticationService
}

func NewAuthenticationHttpController(authenticationService *AuthenticationService) AuthenticationHttpController {
	controller := AuthenticationHttpController{
		Mux:                   chi.NewRouter(),
		authenticationService: authenticationService,
	}

	controller.Post("/authenticate", controller.Authenticate)
	controller.Post("/register", controller.Register)

	return controller
}

type RegisterForm struct {
	Identifier string `form:"identifier"`
	Password   string `form:"password"`
}

func (controller *AuthenticationHttpController) Register(w http.ResponseWriter, req *http.Request) {
	var err error
	form := new(RegisterForm)
	err = json.NewDecoder(req.Body).Decode(form)
	if err != nil {
		common.WriteError(w, http.StatusBadRequest, err)
		return
	}

	err = controller.authenticationService.Register(req.Context(), form.Identifier, form.Password)
	if err != nil {
		common.WriteError(w, http.StatusBadRequest, err)
		return
	}

	common.WriteResponse(w, http.StatusCreated, common.OperationDto{Success: true})
}

type AuthenticateForm struct {
	Identifier string `form:"identifier"`
	Password   string `form:"password"`
}

func (controller *AuthenticationHttpController) Authenticate(w http.ResponseWriter, req *http.Request) {
	var err error
	form := new(AuthenticateForm)

	err = json.NewDecoder(req.Body).Decode(form)
	if err != nil {
		common.WriteError(w, http.StatusBadRequest, err)
		return
	}

	token, err := controller.authenticationService.Authenticate(req.Context(), form.Identifier, form.Password)
	if err != nil {
		common.WriteError(w, http.StatusBadRequest, err)
		return
	}

	common.WriteResponse(w, http.StatusOK, common.AnyDto{"token": token})
}
