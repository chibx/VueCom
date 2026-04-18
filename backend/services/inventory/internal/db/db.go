package db

import (
	"context"
	"errors"

	inventoryModel "github.com/chibx/vuecom/backend/shared/models/db/inventory"
	"gorm.io/gorm"
)

var errDbUnimplemented = errors.New("unimplemented")

type InventoryDB struct {
	db *gorm.DB
}

func NewInventoryDB(db *gorm.DB) *InventoryDB {
	return &InventoryDB{
		db: db,
	}
}

func (c *InventoryDB) CreateInventoryRecords(ctx context.Context, inventoryItems []*inventoryModel.Inventory) error {
	return c.db.WithContext(ctx).Create(inventoryItems).Error
}

func (c *InventoryDB) HasAnyWarehouse(ctx context.Context) (bool, error) {
	var count int64
	err := c.db.WithContext(ctx).Model(&inventoryModel.Warehouse{}).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
