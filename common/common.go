package common

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"log"
	"os"

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

func PanicOnError(err error) {
	if err != nil {
		log.Panicf("unhandled error: %s", err.Error())
	}
}

func LoadPrivateKey(filename string) rsa.PrivateKey {
	privatePem, err := os.ReadFile(filename)
	PanicOnError(err)
	privateBlock, _ := pem.Decode(privatePem)
	print(privateBlock.Bytes)
	privateKey, err := x509.ParsePKCS1PrivateKey(privateBlock.Bytes)
	PanicOnError(err)
	return *privateKey
}
