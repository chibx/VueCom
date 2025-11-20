package main

import (
	model "vuecom/shared/models/db"

	"gorm.io/gorm"
)

func migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		model.BackendUser{},
		model.Product{},
		model.Customer{},
		model.OTP{},
		model.CustomerSession{},
		model.BackendSession{},
	)
}
