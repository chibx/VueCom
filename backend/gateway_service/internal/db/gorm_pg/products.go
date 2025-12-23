package gorm_pg

import (
	"context"

	"vuecom/gateway/internal/types"
	dbModels "vuecom/shared/models/db"

	"gorm.io/gorm"
)

type productRepository struct {
	db *gorm.DB
}

func (p *productRepository) GetProductById(id int, ctx context.Context) (*dbModels.Product, error) {
	return nil, types.ErrDbUnimplemented
}
