//go:build spec

package commands

import (
	"context"
	"go-microservices/resources/inventory"
	"testing"

	"github.com/stretchr/testify/assert"
)

func SetupHarvestCommandTest() (HarvestCommandHandler, func()) {
	inventoryStore := inventory.NewFakeInventoryStore()
	handler := NewHarvestCommandHandler(&inventoryStore)
	return handler, func() {
		inventoryStore.Clear()
	}
}

func TestHarvestCommandHandler(t *testing.T) {
	assert := assert.New(t)
	handler, cleanup := SetupHarvestCommandTest()

	t.Run("it should fail to harvest for an owner that has no inventory", func(t *testing.T) {
		t.Cleanup(cleanup)

		command := HarvestCommand{owner: inventory.OwnerID("who-am-i")}
		err := handler.Handle(context.Background(), command)

		assert.ErrorIs(err, inventory.ErrInventoryNotFound)
		_, err = handler.inventoryStore.Get(context.Background(), command.owner)
		assert.ErrorIs(err, inventory.ErrInventoryNotFound)
	})

	t.Run("it should successfully harvest for an owner that has inventory", func(t *testing.T) {
		t.Cleanup(cleanup)
		owner := inventory.OwnerID("ownerId")
		invent := inventory.NewInventory(owner, 300)
		handler.inventoryStore.Save(context.Background(), invent)

		command := HarvestCommand{owner: owner}
		err := handler.Handle(context.Background(), command)

		assert.NoError(err)
		invent, _ = handler.inventoryStore.Get(context.Background(), command.owner)
		assert.Greater(invent.GetCurrentVolume(), uint(0))
	})

	t.Run("it should fail to harvest for inventory that doesn't have the space", func(t *testing.T) {
		t.Cleanup(cleanup)
		owner := inventory.OwnerID("ownerId")
		invent := inventory.NewInventory(owner, 100)
		handler.inventoryStore.Save(context.Background(), invent)

		command := HarvestCommand{owner: owner}
		err := handler.Handle(context.Background(), command)

		assert.ErrorContains(err, "inventory cannot store item Apple")
		assert.ErrorContains(err, ", because current volume is 0/100")
		invent, _ = handler.inventoryStore.Get(context.Background(), command.owner)
		assert.Equal(invent.GetCurrentVolume(), uint(0))
	})
}
