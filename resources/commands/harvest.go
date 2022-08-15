package commands

import (
	"context"
	"go-microservices/resources/inventory"
	"go-microservices/resources/resource"
)

type HarvestCommand struct {
	owner inventory.OwnerID
}

type HarvestCommandHandler struct {
	inventoryStore inventory.InventoryStore
}

func NewHarvestCommandHandler(inventoryStore inventory.InventoryStore) HarvestCommandHandler {
	return HarvestCommandHandler{inventoryStore: inventoryStore}
}

func (handler *HarvestCommandHandler) Handle(ctx context.Context, command HarvestCommand) error {
	inventory, err := handler.inventoryStore.Get(ctx, command.owner)
	if err != nil {
		return err
	}

	item := resource.Apple.GenerateRandomItem()

	err = inventory.Store(item)
	if err != nil {
		return err
	}

	return handler.inventoryStore.Save(ctx, inventory)
}
