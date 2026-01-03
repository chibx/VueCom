package database

import (
	"context"
	dbModels "vuecom/shared/models/db"
)

type ProductRepository interface {
	CreateProduct(product *dbModels.Product, ctx context.Context) error
	GetProductById(id int, ctx context.Context) (*dbModels.Product, error)
}
