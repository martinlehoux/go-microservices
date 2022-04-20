package main

import (
	"go-microservices/authentication"
	"go-microservices/common"
	"go-microservices/human_resources"
	"log"
	"net/http"
)

func main() {
	privateKey := common.LoadPrivateKey("id_rsa")
	authenticationController := authentication.Bootstrap("/auth", privateKey)
	humanResourcesController := human_resources.Bootstrap("/users", privateKey.PublicKey)

	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		common.WriteResponse(w, http.StatusOK, common.AnyDto{"message": "pong"})
	})

	http.Handle("/users/", humanResourcesController)
	http.Handle("/auth/", authenticationController)

	err := http.ListenAndServe(":3000", nil)
	log.Fatal(err)

}
