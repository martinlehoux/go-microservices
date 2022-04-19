package human_resources

import (
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"go-microservices/common"
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

type UserRegisterForm struct {
	PreferredName string `json:"preferredName"`
}

func (controller *HumanResourcesHttpController) Register(w http.ResponseWriter, req *http.Request) {
	var err error
	token, err := common.ExtractToken(req.Header.Get("Authorization"), controller.publicKey)
	if err != nil {
		common.WriteResponse(w, http.StatusUnauthorized, common.Data{"error": fmt.Sprintf("Unauthorized: %s", err)})
		return
	}
	form := new(UserRegisterForm)
	err = json.NewDecoder(req.Body).Decode(form)
	if err != nil {
		common.WriteResponse(w, http.StatusBadRequest, common.Data{"error": err.Error()})
		return
	}
	err = controller.humanResourcesService.Register(token.Identifier, form.PreferredName)
	if err != nil {
		common.WriteResponse(w, http.StatusBadRequest, common.Data{"error": err.Error()})
		return
	}

	common.WriteResponse(w, http.StatusCreated, common.Data{"success": true})
}

func (controller *HumanResourcesHttpController) GetUsers(w http.ResponseWriter, req *http.Request) {
	users, err := controller.humanResourcesService.GetUsers()
	if err != nil {
		common.WriteResponse(w, http.StatusBadRequest, common.Data{"error": err.Error()})
	}

	common.WriteResponse(w, http.StatusOK, common.Data{"items": users, "total": len(users)})
}
