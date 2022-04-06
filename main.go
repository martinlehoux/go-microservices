package main

import (
	"go-microservices/authentication"
	"go-microservices/user"
	"log"

	"github.com/gofiber/fiber/v2"
)

type UserRegisterForm struct {
	Email string `form:"email"`
}

func main() {
	authenticationService := authentication.Bootstrap()
	userService := user.Bootstrap(authenticationService)

	app := fiber.New()

	app.Get("/ping", func(ctx *fiber.Ctx) error {
		return ctx.SendString("pong")
	})

	app.Get("/users/", func(ctx *fiber.Ctx) error {
		users, err := userService.GetUsers()
		if err != nil {
			return err
		}
		return ctx.JSON(&fiber.Map{
			"total": len(users),
			"items": users,
		})
	})

	app.Post("/users/register", func(ctx *fiber.Ctx) error {
		var err error
		form := new(UserRegisterForm)
		err = ctx.BodyParser(form)
		if err != nil {
			return err
		}
		err = userService.Register(form.Email)
		if err != nil {
			return ctx.Status(400).JSON(&fiber.Map{"error": err.Error()})
		}
		return ctx.Status(201).JSON(&fiber.Map{"success": true})
	})

	log.Fatal(app.Listen(":3000"))

}
