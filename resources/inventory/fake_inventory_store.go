package inventory

import "context"

type FakeInventoryStore struct {
	inventories map[OwnerID]Inventory
}

func NewFakeInventoryStore() FakeInventoryStore {
	return FakeInventoryStore{
		inventories: make(map[OwnerID]Inventory),
	}
}

func (store *FakeInventoryStore) Get(ctx context.Context, owner OwnerID) (Inventory, error) {
	inventory, found := store.inventories[owner]
	if !found {
		return Inventory{}, ErrInventoryNotFound
	}
	return inventory, nil
}

func (store *FakeInventoryStore) Save(ctx context.Context, inventory Inventory) error {
	store.inventories[inventory.owner] = inventory
	return nil
}

func (store *FakeInventoryStore) Clear() {
	store.inventories = make(map[OwnerID]Inventory)
}
