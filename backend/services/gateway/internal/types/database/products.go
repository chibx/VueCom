package database

import (
	"context"

	productModels "github.com/chibx/vuecom/backend/shared/models/db/catalog"
)

type ProductRepository interface {
	CreateProduct(product *productModels.Product, ctx context.Context) error
	GetProductById(id int, ctx context.Context) (*productModels.Product, error)
}
