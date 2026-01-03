package gorm_pg

import (
	"context"

	dbModels "vuecom/shared/models/db"

	"gorm.io/gorm"
)

type productRepository struct {
	db *gorm.DB
}

func (p *productRepository) CreateProduct(product *dbModels.Product, ctx context.Context) error {
	return p.db.WithContext(ctx).Create(product).Error
}

func (p *productRepository) GetProductById(id int, ctx context.Context) (*dbModels.Product, error) {
	product := &dbModels.Product{}

	err := p.db.WithContext(ctx).Where("id = ?", id).First(product).Error
	if err != nil {
		return nil, err
	}

	return product, nil
}
