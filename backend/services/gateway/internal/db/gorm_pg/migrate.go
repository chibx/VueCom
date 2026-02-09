package gorm_pg

import (
	"gorm.io/gorm"
	// model "github.com/chibx/vuecom/backend/shared/models/db"
	// appdata "github.com/chibx/vuecom/backend/shared/models/db/appdata"
	// catalog "github.com/chibx/vuecom/backend/shared/models/db/catalog"
	// inventory "github.com/chibx/vuecom/backend/shared/models/db/inventory"
	// orders "github.com/chibx/vuecom/backend/shared/models/db/orders"
	// users "github.com/chibx/vuecom/backend/shared/models/db/users"
)

func migrate(db *gorm.DB) error {
	return nil // We dont need this
	// var err error

	// err = db.SetupJoinTable(&model.Product{}, "Categories", &model.ProductCategoryValues{})
	// if err != nil {
	// 	return err
	// }
	// err = db.SetupJoinTable(&model.Product{}, "Tags", &model.ProductTags{})
	// if err != nil {
	// 	return err
	// }

	// return db.AutoMigrate(
	// 	// Important Tables
	// 	appdata.AppData{},
	// 	users.Country{},
	// 	users.State{},
	// 	// Backend
	// 	users.BackendUser{},
	// 	users.ApiKey{},
	// 	users.BackendOTP{},
	// 	users.BackendSession{},
	// 	users.BackendUserActivity{},
	// 	users.BackendPasswordResetRequest{},
	// 	// Customer
	// 	users.Customer{},
	// 	users.CustomerOTP{},
	// 	users.CustomerSession{},
	// 	users.CartItem{},
	// 	users.WishlistItem{},
	// 	// Catalog
	// 	catalog.Product{},
	// 	catalog.Attribute{},
	// 	catalog.Category{},
	// 	catalog.Preset{},
	// 	catalog.PresetAttributes{},
	// 	catalog.ProductCategoryValues{},
	// 	catalog.Tag{},
	// 	catalog.ProductTags{},
	// 	catalog.PromoCode{},
	// 	catalog.PromoCodeUsage{},
	// 	// Orders
	// 	orders.Order{},
	// 	orders.OrderItem{},
	// 	orders.OrderReturn{},
	// 	// Inventory
	// 	inventory.Inventory{},
	// 	inventory.Warehouse{},
	// 	inventory.StockMovement{},
	// )
}
