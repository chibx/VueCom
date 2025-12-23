package gorm_pg

import (
	"vuecom/gateway/internal/types"

	"gorm.io/gorm"
)

type GormPGDatabase struct {
	db *gorm.DB
}

func NewGormPGDatabase(db *gorm.DB) *GormPGDatabase {
	return &GormPGDatabase{db: db}
}

func (d *GormPGDatabase) BackendUsers() types.BackendUserRepository {
	return &backendUserRepository{db: d.db}
}

func (d *GormPGDatabase) Customers() types.CustomerRepository {
	return &customerRepository{db: d.db}
}

func (d *GormPGDatabase) Categories() types.CategoryRepository {
	return &categoryRepository{db: d.db}
}

func (d *GormPGDatabase) Products() types.ProductRepository {
	return &productRepository{db: d.db}
}

func (d *GormPGDatabase) Orders() types.OrderRepository {
	return &orderRepository{db: d.db}
}

func (d *GormPGDatabase) Inventory() types.InventoryRepository {
	return &inventoryRepository{db: d.db}
}

func (d *GormPGDatabase) AppData() types.AppDataRepository {
	return &appdataRepository{db: d.db}
}
