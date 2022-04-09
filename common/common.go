package common

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type ID = uuid.UUID

func CreateID() ID {
	return uuid.New()
}

func SendError(ctx *fiber.Ctx, err error) error {
	return ctx.Status(400).JSON(&fiber.Map{"error": err.Error()})
}
