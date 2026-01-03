package database

import (
	"context"
	dbModels "vuecom/shared/models/db"
)

type CategoryRepository interface {
	GetCategoryById(id int, ctx context.Context) (*dbModels.Category, error)
}
