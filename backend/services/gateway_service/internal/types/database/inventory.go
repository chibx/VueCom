package database

import (
	"context"
	invModels "vuecom/shared/models/db/inventory"
)

type InventoryRepository interface {
	GetInventoryById(id int, ctx context.Context) (*invModels.Inventory, error)
}
