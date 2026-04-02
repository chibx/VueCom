package db

import (
	"context"

	inventoryModel "github.com/chibx/vuecom/backend/shared/models/db/inventory"
)

func (c *InventoryDB) GetInventoryById(id int, ctx context.Context) (*inventoryModel.Inventory, error) {
	return nil, errDbUnimplemented
}
