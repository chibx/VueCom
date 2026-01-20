package database

import (
	"context"

	invModels "github.com/chibx/vuecom/backend/shared/models/db/inventory"
)

type InventoryRepository interface {
	GetInventoryById(id int, ctx context.Context) (*invModels.Inventory, error)
}
