package gorm_pg

import (
	"context"

	catalogModels "vuecom/shared/models/db/catalog"

	"gorm.io/gorm"
)

type categoryRepository struct {
	db *gorm.DB
}

func (c *categoryRepository) GetCategoryById(id int, ctx context.Context) (*catalogModels.Category, error) {
	return nil, errDbUnimplemented
}
