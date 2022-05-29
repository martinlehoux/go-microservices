package main

import (
	"context"
	"fmt"
	"go-microservices/authentication"
	"go-microservices/common"
	"go-microservices/human_resources"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

func main() {
	logger := common.NewLogrusLogger()
	start := time.Now()
	logger.Info(context.Background(), "booting up services")
	privateKey := common.LoadPrivateKey("id_rsa")
	authenticationController := authentication.Bootstrap(&logger, privateKey)
	humanResourcesController := human_resources.Bootstrap(&logger, privateKey.PublicKey)

	router := chi.NewRouter()

	router.Use(common.CommonMiddlewareConstructor(&logger))

	router.Get("/ping", func(w http.ResponseWriter, req *http.Request) {
		common.WriteResponse(w, http.StatusOK, common.AnyDto{"message": "pong"})
	})

	router.Mount("/hr", humanResourcesController)
	router.Mount("/auth", authenticationController)
	logger.Info(context.Background(), "services boot up in %d ms", time.Since(start).Milliseconds())

	const port = 3000
	logger.Info(context.Background(), "listening on http://localhost:%d", port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), router)
	log.Fatal(err.Error())
}
