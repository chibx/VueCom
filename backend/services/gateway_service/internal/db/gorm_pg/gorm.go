package gorm_pg

import (
	"errors"
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

var errDbUnimplemented = errors.New("unimplemented")

type GormPGDatabase struct {
	db *gorm.DB
}

func NewGormPGDatabase(db *gorm.DB) *GormPGDatabase {
	return &GormPGDatabase{db: db}
}

func (d *GormPGDatabase) BackendUsers() database.BackendUserRepository {
	if backendR == nil {
		backendR = &backendUserRepository{db: d.db}
	}
	return backendR
}

func (d *GormPGDatabase) Customers() database.CustomerRepository {
	if customerR == nil {
		customerR = &customerRepository{db: d.db}
	}
	return customerR
}

func (d *GormPGDatabase) Categories() database.CategoryRepository {
	if categoryR == nil {
		categoryR = &categoryRepository{db: d.db}
	}
	return categoryR
}

func (d *GormPGDatabase) Products() database.ProductRepository {
	if productR == nil {
		productR = &productRepository{db: d.db}
	}
	return productR
}

func (d *GormPGDatabase) Orders() database.OrderRepository {
	if ordersR == nil {
		ordersR = &orderRepository{db: d.db}
	}
	return ordersR
}

func (d *GormPGDatabase) Inventory() database.InventoryRepository {
	if inventoryR == nil {
		inventoryR = &inventoryRepository{db: d.db}
	}
	return inventoryR
}

func (d *GormPGDatabase) AppData() database.AppDataRepository {
	if appdataR == nil {
		appdataR = &appdataRepository{db: d.db}
	}
	return appdataR
}

func (d *GormPGDatabase) Migrate() error {
	return migrate(d.db)
}
