package database

import (
	"context"

	catModels "github.com/chibx/vuecom/backend/shared/models/db/catalog"
)

type CategoryRepository interface {
	GetCategoryById(id int, ctx context.Context) (*catModels.Category, error)
}
