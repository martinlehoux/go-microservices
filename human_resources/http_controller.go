package human_resources

import (
	"go-microservices/common"

	"github.com/gofiber/fiber/v2"
)

type UserRegisterForm struct {
	Email string `form:"email"`
}

func BootstrapHttpController(router fiber.Router, userService *UserService) {
	router.Get("/", func(ctx *fiber.Ctx) error {
		users, err := userService.GetUsers()
		if err != nil {
			return err
		}
		return ctx.JSON(&fiber.Map{
			"total": len(users),
			"items": users,
		})
	})

	router.Post("/register", func(ctx *fiber.Ctx) error {
		var err error
		form := new(UserRegisterForm)
		err = ctx.BodyParser(form)
		if err != nil {
			return common.SendError(ctx, err)
		}
		err = userService.Register(form.Email)
		if err != nil {
			return common.SendError(ctx, err)
		}
		return ctx.Status(201).JSON(&fiber.Map{"success": true})
	})
}
