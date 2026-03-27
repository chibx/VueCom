package catalog

import (
	"time"

	"github.com/chibx/vuecom/backend/shared/models/db/catalog"
)

type CreateProductReq struct {
	P             catalog.Product
	Name          string     `json:"name"`
	SKU           string     `json:"sku"`
	BasePrice     float64    `json:"base_price"`
	SalePrice     float64    `json:"sale_price"`
	DiscountStart *time.Time `json:"discount_start"`
	DiscountEnd   *time.Time `json:"discount_end"`
	Enabled       bool       `json:"enabled"`
	// Warranty   string    `json:"warranty"`
	ShortDescription string   `json:"short_description"`
	FullDescription  string   `json:"full_description"`
	Quantity         int      `json:"quantity"`
	Slug             string   `json:"slug"`
	Country          uint     `json:"country"`
	Weight           *float64 `json:"weight"`
	BrandId          int      `json:"brand_id"`
	Color            uint     `json:"color"`
	Medias           []uint   `json:"medias"`
	MetaTitle        *string  `json:"meta_title,omitempty"`
	MetaDescription  *string  `json:"meta_description,omitempty"`
	SearchKeywords   *string  `json:"search_keywords"`
	RelatedProducts  []string `json:"related_products"`
	UpSellProducts   []string `json:"upsell"`
	ParentID         *uint    `json:"parent_id"`
	PresetID         *uint    `json:"preset_id"`
}
