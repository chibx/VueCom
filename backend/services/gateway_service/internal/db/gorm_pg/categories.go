package gorm_pg

import (
	"context"

	dbModels "vuecom/shared/models/db"

	"gorm.io/gorm"
)

type categoryRepository struct {
	db *gorm.DB
}

func (c *categoryRepository) GetCategoryById(id int, ctx context.Context) (*dbModels.Category, error) {
	return nil, errDbUnimplemented
}
