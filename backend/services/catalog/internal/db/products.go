package db

import (
	"context"

	serverErrors "github.com/chibx/vuecom/backend/shared/errors/server"
	catModels "github.com/chibx/vuecom/backend/shared/models/db/catalog"

	"gorm.io/gorm"
)

func (c *CatalogDB) CreateProduct(ctx context.Context, product *catModels.Product) error {
	return c.db.WithContext(ctx).Create(product).Error
}

func (c *CatalogDB) GetProductById(ctx context.Context, id int) (*catModels.Product, error) {
	product := &catModels.Product{}

	err := c.db.WithContext(ctx).Where("id = ?", id).First(product).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, serverErrors.ErrDBRecordNotFound
		}
		return nil, err
	}

	return product, nil
}
