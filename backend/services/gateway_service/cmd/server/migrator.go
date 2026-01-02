package main

import (
	model "vuecom/shared/models/db"

	"gorm.io/gorm"
)

func migrate(db *gorm.DB) error {
	var err error
	// I would add this to a docker init script for PostgreSQL
	err = db.Exec(`CREATE SCHEMA IF NOT EXISTS backend;
	CREATE SCHEMA IF NOT EXISTS customer;
	CREATE SCHEMA IF NOT EXISTS catalog;
	CREATE SCHEMA IF NOT EXISTS inventory;`).Error
	if err != nil {
		return err
	}

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
		model.AppData{},
		model.Country{},
		model.State{},
		// Backend
		model.BackendUser{},
		model.ApiKey{},
		model.BackendOTP{},
		model.BackendSession{},
		model.BackendUserActivity{},
		model.BackendPasswordResetRequest{},
		// Catalog
		model.Product{},
		model.Attribute{},
		model.Category{},
		model.Preset{},
		model.PresetAttributes{},
		model.ProductCategoryValues{},
		model.Tag{},
		model.ProductTags{},
		model.PromoCode{},
		model.PromoCodeUsage{},
		model.Order{},
		model.OrderItem{},
		model.OrderReturn{},
		// Customer
		model.Customer{},
		model.CustomerOTP{},
		model.CustomerSession{},
		model.CartItem{},
		model.WishlistItem{},

		// Inventory
		model.Inventory{},
		model.Warehouse{},
		model.StockMovement{},
	)
}
