package db

import (
	"errors"

	"gorm.io/gorm"
)

var errDbUnimplemented = errors.New("unimplemented")

type InventoryDB struct {
	db *gorm.DB
}

func NewInventoryDB(db *gorm.DB) *InventoryDB {
	return &InventoryDB{
		db: db,
	}
}
