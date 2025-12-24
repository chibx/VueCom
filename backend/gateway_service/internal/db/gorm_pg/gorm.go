package gorm_pg

import (
	"vuecom/gateway/internal/types"

	"gorm.io/gorm"
)

var backendR *backendUserRepository
var customerR *customerRepository
var productR *productRepository
var categoryR *categoryRepository
var ordersR *orderRepository
var inventoryR *inventoryRepository
var appdataR *appdataRepository

type gormPGDatabase struct {
	db *gorm.DB
}

func NewGormPGDatabase(db *gorm.DB) *gormPGDatabase {
	return &gormPGDatabase{db: db}
}

func (d *gormPGDatabase) BackendUsers() types.BackendUserRepository {
	if backendR == nil {
		backendR = &backendUserRepository{db: d.db}
	}
	return backendR
}

func (d *gormPGDatabase) Customers() types.CustomerRepository {
	if customerR == nil {
		customerR = &customerRepository{db: d.db}
	}
	return customerR
}

func (d *gormPGDatabase) Categories() types.CategoryRepository {
	if categoryR == nil {
		categoryR = &categoryRepository{db: d.db}
	}
	return categoryR
}

func (d *gormPGDatabase) Products() types.ProductRepository {
	if productR == nil {
		productR = &productRepository{db: d.db}
	}
	return productR
}

func (d *gormPGDatabase) Orders() types.OrderRepository {
	if ordersR == nil {
		ordersR = &orderRepository{db: d.db}
	}
	return ordersR
}

func (d *gormPGDatabase) Inventory() types.InventoryRepository {
	if inventoryR == nil {
		inventoryR = &inventoryRepository{db: d.db}
	}
	return inventoryR
}

func (d *gormPGDatabase) AppData() types.AppDataRepository {
	if appdataR == nil {
		appdataR = &appdataRepository{db: d.db}
	}
	return appdataR
}

func (d *gormPGDatabase) Migrate() error {
	return migrate(d.db)
}
