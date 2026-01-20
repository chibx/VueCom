package gorm_pg

import (
	"gorm.io/gorm"

	model "vuecom/shared/models/db"
	appdata "vuecom/shared/models/db/appdata"
	catalog "vuecom/shared/models/db/catalog"
	inventory "vuecom/shared/models/db/inventory"
	orders "vuecom/shared/models/db/orders"
	users "vuecom/shared/models/db/users"
)

func migrate(db *gorm.DB) error {
	var err error

	err = db.SetupJoinTable(&model.Product{}, "Categories", &model.ProductCategoryValues{})
	if err != nil {
		return err
	}
	err = db.SetupJoinTable(&model.Product{}, "Tags", &model.ProductTags{})
	if err != nil {
		return err
	}

	return db.AutoMigrate(
		// Important Tables
		appdata.AppData{},
		users.Country{},
		users.State{},
		// Backend
		users.BackendUser{},
		users.ApiKey{},
		users.BackendOTP{},
		users.BackendSession{},
		users.BackendUserActivity{},
		users.BackendPasswordResetRequest{},
		// Customer
		users.Customer{},
		users.CustomerOTP{},
		users.CustomerSession{},
		users.CartItem{},
		users.WishlistItem{},
		// Catalog
		catalog.Product{},
		catalog.Attribute{},
		catalog.Category{},
		catalog.Preset{},
		catalog.PresetAttributes{},
		catalog.ProductCategoryValues{},
		catalog.Tag{},
		catalog.ProductTags{},
		catalog.PromoCode{},
		catalog.PromoCodeUsage{},
		// Orders
		orders.Order{},
		orders.OrderItem{},
		orders.OrderReturn{},
		// Inventory
		inventory.Inventory{},
		inventory.Warehouse{},
		inventory.StockMovement{},
	)
}
