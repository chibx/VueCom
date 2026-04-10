package db

import (
	"context"
	"fmt"
	"strings"

	"github.com/chibx/vuecom/backend/shared/constants"
	serverErrors "github.com/chibx/vuecom/backend/shared/errors/server"
	"github.com/chibx/vuecom/backend/shared/models/db/catalog"
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
		fmt.Fprintf(&sql, "(%d, %d)",
			productId,
			id,
		)

		if idx < len(categories)-1 {
			sql.WriteString(",")
		}
	}
	sql.WriteString(";")

	return c.db.WithContext(ctx).Exec(sql.String()).Error
}

func (c *CatalogDB) CreateProductRelation(ctx context.Context, productId uint32, related []uint32, upsell []uint32, crossSell []uint32) error {
	// [RAW SQL] To avoid the allocation of an array/slice of catalog.ProductRelation struct
	// Is safe (supposedly) because I am inserting integers, not strings
	if len(related) == 0 && len(upsell) == 0 && len(crossSell) == 0 {
		return nil
	}

	var relatedSql strings.Builder
	var sortOrder = 0
	_ = catalog.ProductRelation{}

	relatedSql.WriteString("INSERT INTO product_relations (source_product_id, target_product_id, relation_type, sort_order) VALUES ")
	for _, related_id := range related {
		fmt.Fprintf(&relatedSql, "(%d, %d, %d, %d),",
			productId,
			related_id,
			uint(constants.RelationRelated),
			sortOrder,
		)
		sortOrder += 1
	}

	// Upsell
	sortOrder = 0
	for _, upsell_id := range upsell {
		fmt.Fprintf(&relatedSql, "(%d, %d, %d, %d),",
			productId,
			upsell_id,
			uint(constants.RelationUpsell),
			sortOrder,
		)
		sortOrder += 1
	}

	// Cross-Sell
	sortOrder = 0
	if len(crossSell) > 0 {
		for _, crossSell_id := range crossSell {
			fmt.Fprintf(&relatedSql, "(%d, %d, %d, %d),",
				productId,
				crossSell_id,
				uint(constants.RelationCrossSell),
				sortOrder,
			)
			sortOrder += 1
		}
	}
	var tmp = relatedSql.String()
	var finalSQL = tmp[:len(tmp)-1] + ";" // Remove the trailing comma
	return c.db.WithContext(ctx).Exec(finalSQL).Error
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
