package gorm_pg

import (
	"context"

	"vuecom/gateway/internal/types"
	dbModels "vuecom/shared/models/db"

	"gorm.io/gorm"
)

type inventoryRepository struct {
	db *gorm.DB
}

func (i *inventoryRepository) GetInventoryById(id int, ctx context.Context) (*dbModels.Inventory, error) {
	return nil, types.ErrDbUnimplemented
}
