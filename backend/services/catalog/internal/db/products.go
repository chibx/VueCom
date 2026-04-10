package db

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	serverErrors "github.com/chibx/vuecom/backend/shared/errors/server"
	catModels "github.com/chibx/vuecom/backend/shared/models/db/catalog"

	"gorm.io/gorm"
)

func (c *CatalogDB) CreateProduct(ctx context.Context, product *catModels.Product) error {
	return c.db.WithContext(ctx).Create(product).Error
}

func (c *CatalogDB) CreateProductToCategory(ctx context.Context, productId uint32, categories []uint32) error {
	// [RAW SQL] To avoid the allocation of an array/slice of catalog.ProductCategoryValues struct
	// Is safe (supposedly) because I am inserting integers, not strings
	if len(categories) == 0 {
		return nil
	}
	var sql strings.Builder
	sql.WriteString(`INSERT INTO product_category_values (product_id, category_id) VALUES `)
	for idx, id := range categories {
		fmt.Fprintf(&sql, "(%s, %s)",
			strconv.FormatInt(int64(productId), 10),
			strconv.FormatInt(int64(id), 10),
		)

		if idx < len(categories)-1 {
			sql.WriteString(",")
		}
	}
	sql.WriteString(";")

	return c.db.WithContext(ctx).Exec(sql.String()).Error
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
