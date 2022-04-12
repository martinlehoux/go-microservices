package main

import (
	"go-microservices/authentication"
	"go-microservices/common"
	"go-microservices/human_resources"
	"log"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()
	authenticationService := authentication.Bootstrap(common.LoadPrivateKey("id_rsa"))
	authentication.BootstrapHttpController(app.Group("/auth"), authenticationService)
	userService := human_resources.Bootstrap()
	human_resources.BootstrapHttpController(app.Group("/users"), userService)

	app.Get("/ping", func(ctx *fiber.Ctx) error {
		return ctx.SendString("pong")
	})

	log.Fatal(app.Listen(":3000"))

}
