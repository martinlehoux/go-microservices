package authentication

import (
	"go-microservices/common"

	"github.com/gofiber/fiber/v2"
)

type RegisterForm struct {
	Identifier string `form:"identifier"`
	Password   string `form:"password"`
}

type AuthenticateForm struct {
	Identifier string `form:"identifier"`
	Password   string `form:"password"`
}

func BootstrapHttpController(router fiber.Router, authenticationService *AuthenticationService) {
	router.Post("/register", func(ctx *fiber.Ctx) error {
		var err error
		form := new(RegisterForm)
		err = ctx.BodyParser(form)
		if err != nil {
			return err
		}
		err = authenticationService.Register(form.Identifier, form.Password)
		if err != nil {
			return common.SendError(ctx, err)
		}
		return ctx.Status(201).JSON(&fiber.Map{"success": true})
	})

	router.Post("/authenticate", func(ctx *fiber.Ctx) error {
		var err error
		form := new(AuthenticateForm)
		err = ctx.BodyParser(form)
		if err != nil {
			return common.SendError(ctx, err)
		}
		token, signature, err := authenticationService.Authenticate(form.Identifier, form.Password)
		if err != nil {
			return common.SendError(ctx, err)
		}
		return ctx.JSON(&fiber.Map{"token": token, "signature": signature})
	})
}
