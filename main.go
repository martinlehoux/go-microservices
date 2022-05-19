package main

import (
	"fmt"
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

	http.HandleFunc("/ping", func(w http.ResponseWriter, req *http.Request) {
		common.CommonMiddleware(w, req)
		common.WriteResponse(w, http.StatusOK, common.AnyDto{"message": "pong"})
	})

	http.Handle("/users/", humanResourcesController)
	http.Handle("/auth/", authenticationController)

	const port = 3000
	log.Printf("listening on http://localhost:%d", port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	log.Fatal(err.Error())

}
