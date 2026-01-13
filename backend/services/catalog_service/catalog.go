package catalog_service

import (
	catalog_db "vuecom/catalog/internal/db"

	"gorm.io/gorm"
)

type CatalogService struct {
	repo *catalog_db.CatalogDB
}

func NewCatalogService(db *gorm.DB) *CatalogService {
	return &CatalogService{
		repo: catalog_db.NewCatalogDB(db),
	}
}
