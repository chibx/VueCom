package gorm_pg

import (
	"context"

	catModels "github.com/chibx/vuecom/backend/shared/models/db/catalog"

	"gorm.io/gorm"
)

type productRepository struct {
	db *gorm.DB
}

func (p *productRepository) CreateProduct(product *catModels.Product, ctx context.Context) error {
	return p.db.WithContext(ctx).Create(product).Error
}

func (p *productRepository) GetProductById(id int, ctx context.Context) (*catModels.Product, error) {
	product := &catModels.Product{}

	err := p.db.WithContext(ctx).Where("id = ?", id).First(product).Error
	if err != nil {
		return nil, err
	}

	return product, nil
}
