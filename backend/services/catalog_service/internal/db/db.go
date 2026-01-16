package db

import (
	"errors"

	"gorm.io/gorm"
)

var errDbUnimplemented = errors.New("unimplemented")

type CatalogDB struct {
	db *gorm.DB
}

func NewCatalogDB(db *gorm.DB) *CatalogDB {
	return &CatalogDB{
		db: db,
	}
}
