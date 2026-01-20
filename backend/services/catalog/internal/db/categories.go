package db

import (
	"context"
	catalogModel "vuecom/shared/models/db/catalog"
)

func (c *CatalogDB) GetCategoryById(id int, ctx context.Context) (*catalogModel.Category, error) {
	return nil, errDbUnimplemented
}
