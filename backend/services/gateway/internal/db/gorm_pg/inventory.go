package gorm_pg

import (
	"context"

	invModels "github.com/chibx/vuecom/backend/shared/models/db/inventory"

	"gorm.io/gorm"
)

type inventoryRepository struct {
	db *gorm.DB
}

func (i *inventoryRepository) GetInventoryById(id int, ctx context.Context) (*invModels.Inventory, error) {
	return nil, errDbUnimplemented
}
