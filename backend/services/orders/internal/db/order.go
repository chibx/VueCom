package db

import (
	"context"

	catalogModel "github.com/chibx/vuecom/backend/shared/models/db/catalog"
)

func (c *OrderDB) GetCategoryById(id int, ctx context.Context) (*catalogModel.Category, error) {
	return nil, errDbUnimplemented
}
