package human_resources

import (
	"crypto/rsa"
	"fmt"
	"go-microservices/common"

	"github.com/gofiber/fiber/v2"
)

type UserRegisterForm struct {
	PreferredName string `json:"preferredName"`
}

func BootstrapHttpController(router fiber.Router, humanResourcesService *HumanResourcesService, publicKey rsa.PublicKey) {
	router.Get("/", func(ctx *fiber.Ctx) error {
		users, err := humanResourcesService.GetUsers()
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
		token, err := common.ExtractToken(ctx.Get("Authorization"), publicKey)
		if err != nil {
			return ctx.Status(401).JSON(&fiber.Map{"error": fmt.Sprintf("Unauthorized: %s", err)})
		}
		form := new(UserRegisterForm)
		err = ctx.BodyParser(form)
		if err != nil {
			return common.SendError(ctx, err)
		}
		err = humanResourcesService.Register(token.Identifier, form.PreferredName)
		if err != nil {
			return common.SendError(ctx, err)
		}
		return ctx.Status(201).JSON(&fiber.Map{"success": true})
	})
}
