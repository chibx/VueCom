package gorm_pg

import (
	"context"
	"errors"

	"vuecom/gateway/internal/types"
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
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, types.ErrDbNil
		}
		return nil, err
	}

	return product, nil
}
