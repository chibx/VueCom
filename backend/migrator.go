package main

import (
	model "vuecom/models/db"

	"gorm.io/gorm"
)

func migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		model.User{},
		model.Product{},
		model.Customer{},
		model.OTP{},
		model.Session{},
	)
}
