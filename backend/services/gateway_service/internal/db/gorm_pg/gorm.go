package gorm_pg

import (
	"vuecom/gateway/internal/types/database"

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

func (d *gormPGDatabase) BackendUsers() database.BackendUserRepository {
	if backendR == nil {
		backendR = &backendUserRepository{db: d.db}
	}
	return backendR
}

func (d *gormPGDatabase) Customers() database.CustomerRepository {
	if customerR == nil {
		customerR = &customerRepository{db: d.db}
	}
	return customerR
}

func (d *gormPGDatabase) Categories() database.CategoryRepository {
	if categoryR == nil {
		categoryR = &categoryRepository{db: d.db}
	}
	return categoryR
}

func (d *gormPGDatabase) Products() database.ProductRepository {
	if productR == nil {
		productR = &productRepository{db: d.db}
	}
	return productR
}

func (d *gormPGDatabase) Orders() database.OrderRepository {
	if ordersR == nil {
		ordersR = &orderRepository{db: d.db}
	}
	return ordersR
}

func (d *gormPGDatabase) Inventory() database.InventoryRepository {
	if inventoryR == nil {
		inventoryR = &inventoryRepository{db: d.db}
	}
	return inventoryR
}

func (d *gormPGDatabase) AppData() database.AppDataRepository {
	if appdataR == nil {
		appdataR = &appdataRepository{db: d.db}
	}
	return appdataR
}

func (d *gormPGDatabase) Migrate() error {
	return migrate(d.db)
}
