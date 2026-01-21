package gorm_pg

import (
	"context"

	catModels "github.com/chibx/vuecom/backend/shared/models/db/catalog"

	"gorm.io/gorm"
)

type productRepository struct {
	db *gorm.DB
}

func (p *productRepository) CreateProduct(ctx context.Context, product *catModels.Product) error {
	return p.db.WithContext(ctx).Create(product).Error
}

func (p *productRepository) GetProductById(ctx context.Context, id int) (*catModels.Product, error) {
	product := &catModels.Product{}

	err := p.db.WithContext(ctx).Where("id = ?", id).First(product).Error
	if err != nil {
		return nil, err
	}

	return product, nil
}
