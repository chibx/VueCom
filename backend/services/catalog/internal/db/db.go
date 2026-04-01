package db

import (
	"errors"

	"gorm.io/gorm"
)

var errDbUnimplemented = errors.New("unimplemented")

type CatalogDB struct {
	db *gorm.DB
}

func (c *CatalogDB) RunInTx(fn func(c *CatalogDB) error) error {
	return c.db.Transaction(func(tx *gorm.DB) error {
		// Create a NEW store implementation that wraps the transaction instance
		txStore := &CatalogDB{db: tx}

		// Execute the user's function with the transaction-bound store
		return fn(txStore)
	})
}

func NewCatalogDB(db *gorm.DB) *CatalogDB {
	return &CatalogDB{
		db: db,
	}
}
