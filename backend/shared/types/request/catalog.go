package request

import (
	"time"
)

type ProductVisibility int

const (
	VisibilityCatalog ProductVisibility = 0
	VisibilitySearch  ProductVisibility = 1
	VisibilityBoth    ProductVisibility = 2
)

type CreateProductReq struct {
	Name             string         `json:"name" validate:"required"`
	SKU              string         `json:"sku" validate:"required"`
	BasePrice        float64        `json:"base_price" validate:"required"`
	SalePrice        float64        `json:"sale_price" validate:"required"`
	InStock          bool           `json:"in_stock" validate:"required"`
	DiscountStart    *time.Time     `json:"discount_start"`
	DiscountEnd      *time.Time     `json:"discount_end"`
	Categories       []int32        `json:"categories"`
	IsNew            bool           `json:"is_new"`
	NewFrom          *time.Time     `json:"new_from"`
	NewTo            *time.Time     `json:"new_to"`
	Enabled          bool           `json:"enabled"`
	ShortDescription string         `json:"short_description"`
	FullDescription  string         `json:"full_description"`
	Quantity         int32          `json:"quantity" validate:"gte=0"`
	Slug             string         `json:"slug" validate:"required"`
	CountryOfManf    uint32         `json:"country_of_manufacture" validate:"required"`
	Weight           *float64       `json:"weight"`
	BrandId          int32          `json:"brand_id"`
	ColorId          uint32         `json:"color_id"`
	Medias           []uint32       `json:"medias"`
	MetaTitle        string         `json:"meta_title,omitempty" validate:"required"`
	MetaDescription  string         `json:"meta_description,omitempty" validate:"required"`
	SearchKeywords   *string        `json:"search_keywords"`
	RelatedProducts  []uint32       `json:"related_products"`
	UpSellProducts   []uint32       `json:"upsell"`
	CrossSell        []uint32       `json:"cross_sell"`
	PresetID         *uint32        `json:"preset_id"`
	PresetValues     map[string]any `json:"preset_values"`
}
