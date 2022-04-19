package authentication

import (
	"encoding/json"
	"go-microservices/common"
	"net/http"
	"strings"
)

type AuthenticationHttpController struct {
	authenticationService *AuthenticationService
	rootPath              string
}

func (controller *AuthenticationHttpController) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	path := strings.TrimPrefix(req.URL.Path, controller.rootPath)
	if path == "/register" && req.Method == "POST" {
		controller.Register(w, req)
	}
	if path == "/authenticate" && req.Method == "POST" {
		controller.Authenticate(w, req)
	}
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
		common.WriteResponse(w, http.StatusBadRequest, common.Data{"error": err.Error()})
		return
	}
	err = controller.authenticationService.Register(form.Identifier, form.Password)
	if err != nil {
		common.WriteResponse(w, http.StatusBadRequest, common.Data{"error": err.Error()})
		return
	}
	common.WriteResponse(w, http.StatusCreated, common.Data{"success": true})
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
		common.WriteResponse(w, http.StatusBadRequest, common.Data{"error": err.Error()})
		return
	}
	token, err := controller.authenticationService.Authenticate(form.Identifier, form.Password)
	if err != nil {
		common.WriteResponse(w, http.StatusBadRequest, common.Data{"error": err.Error()})
		return
	}
	common.WriteResponse(w, http.StatusOK, common.Data{"token": token})
}
