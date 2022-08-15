package resources

import (
	"crypto/rsa"
	"go-microservices/common"
)

func Bootstrap(logger common.Logger, commandBus) *ResourcesHttpController {
	inventoryStore := NewFakeInventoryStore()
	service := NewHumanResourcesService(&userStore, &groupStore, logger)
	controller := NewHumanResourcesHttpController(&service, publicKey)
	return &controller
}
