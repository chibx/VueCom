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
	Categories       []int          `json:"categories"`
	IsNew            bool           `json:"is_new"`
	NewFrom          time.Time      `json:"new_from"`
	NewTo            time.Time      `json:"new_to"`
	Enabled          bool           `json:"enabled"`
	ShortDescription string         `json:"short_description"`
	FullDescription  string         `json:"full_description"`
	Quantity         int            `json:"quantity"`
	Slug             string         `json:"slug" validate:"required"`
	CountryOfManf    uint           `json:"country_of_manufacture" validate:"required"`
	Weight           *float64       `json:"weight"`
	BrandId          int            `json:"brand_id"`
	ColorId          uint           `json:"color_id"`
	Medias           []uint         `json:"medias"`
	MetaTitle        string         `json:"meta_title,omitempty" validate:"required"`
	MetaDescription  string         `json:"meta_description,omitempty" validate:"required"`
	SearchKeywords   *string        `json:"search_keywords"`
	RelatedProducts  []uint         `json:"related_products"`
	UpSellProducts   []uint         `json:"upsell"`
	CrossSell        []uint         `json:"cross_sell"`
	PresetID         *uint          `json:"preset_id"`
	PresetValues     map[string]any `json:"preset_values"`
}
