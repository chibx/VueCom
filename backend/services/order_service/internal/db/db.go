package db

import (
	"errors"

	"gorm.io/gorm"
)

var errDbUnimplemented = errors.New("unimplemented")

type OrderDB struct {
	db *gorm.DB
}

func NewOrderDB(db *gorm.DB) *OrderDB {
	return &OrderDB{
		db: db,
	}
}
