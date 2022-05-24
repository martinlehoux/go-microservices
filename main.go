package main

import (
	"context"
	"fmt"
	"go-microservices/authentication"
	"go-microservices/common"
	"go-microservices/human_resources"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func main() {
	privateKey := common.LoadPrivateKey("id_rsa")
	logger := common.NewLogrusLogger()
	authenticationController := authentication.Bootstrap(&logger, privateKey)
	humanResourcesController := human_resources.Bootstrap(&logger, privateKey.PublicKey)

	router := chi.NewRouter()

	router.Use(common.CommonMiddlewareConstructor(&logger))

	router.Get("/ping", func(w http.ResponseWriter, req *http.Request) {
		common.WriteResponse(w, http.StatusOK, common.AnyDto{"message": "pong"})
	})

	router.Mount("/hr", humanResourcesController)
	router.Mount("/auth", authenticationController)

	const port = 3000
	logger.Info(context.Background(), "listening on http://localhost:%d", port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), router)
	log.Fatal(err.Error())
}
