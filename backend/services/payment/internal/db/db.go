package db

import (
	"errors"

	"gorm.io/gorm"
)

var errDbUnimplemented = errors.New("unimplemented")

type PaymentDB struct {
	db *gorm.DB
}

func NewPaymentDB(db *gorm.DB) *PaymentDB {
	return &PaymentDB{
		db: db,
	}
}
