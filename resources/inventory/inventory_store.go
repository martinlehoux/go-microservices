package inventory

import (
	"context"
	"errors"
)

type InventoryStore interface {
	Clear()
	Get(ctx context.Context, owner OwnerID) (Inventory, error)
	Save(ctx context.Context, inventory Inventory) error
}

var ErrInventoryNotFound = errors.New("inventory not found")
