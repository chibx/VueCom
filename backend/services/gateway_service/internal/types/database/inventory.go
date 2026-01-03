package database

import (
	"context"
	dbModels "vuecom/shared/models/db"
)

type InventoryRepository interface {
	GetInventoryById(id int, ctx context.Context) (*dbModels.Inventory, error)
}
