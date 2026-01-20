package database

import (
	"context"
	catModels "vuecom/shared/models/db/catalog"
)

type CategoryRepository interface {
	GetCategoryById(id int, ctx context.Context) (*catModels.Category, error)
}
