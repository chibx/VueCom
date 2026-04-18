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

func (c *InventoryDB) ListWarehouses(ctx context.Context) ([]*inventoryModel.Warehouse, error) {
	var warehouses []*inventoryModel.Warehouse
	err := c.db.WithContext(ctx).Find(&warehouses).Error
	return warehouses, err
}

func (c *InventoryDB) CreateWarehouse(ctx context.Context, warehouse *inventoryModel.Warehouse) error {
	return c.db.WithContext(ctx).Create(warehouse).Error
}

func (c *InventoryDB) DeleteWarehouse(ctx context.Context, ids []uint32) error {
	return c.db.WithContext(ctx).Delete(&inventoryModel.Warehouse{}, ids).Error
}

func (c *InventoryDB) CreateStockMovement(ctx context.Context, movement *inventoryModel.StockMovement) error {
	return c.db.WithContext(ctx).Create(movement).Error
}

func (c *InventoryDB) ListStockMovements(ctx context.Context, warehouseId uint, sku string) ([]*inventoryModel.StockMovement, error) {
	var movements []*inventoryModel.StockMovement
	query := c.db.WithContext(ctx)
	if warehouseId > 0 {
		query = query.Where("warehouse_id = ?", warehouseId)
	}
	if sku != "" {
		query = query.Where("sku = ?", sku)
	}
	err := query.Find(&movements).Error
	return movements, err
}
