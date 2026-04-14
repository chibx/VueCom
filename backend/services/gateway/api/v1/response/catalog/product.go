package catalog

import (
	"time"
)

type GetProductResp struct {
	ID               uint32     `json:"id" redis:"id"`
	Name             string     `json:"name" redis:"name"`
	SKU              string     `json:"sku" redis:"sku"`
	BasePrice        float64    `json:"base_price" redis:"base_price"`
	SalePrice        float64    `json:"sale_price" redis:"sale_price"`
	InStock          bool       `json:"in_stock" redis:"in_stock"`
	DiscountStart    *time.Time `json:"discount_start" redis:"discount_start"`
	DiscountEnd      *time.Time `json:"discount_end" redis:"discount_end"`
	Categories       []int32    `json:"categories" redis:"categories"`
	IsNew            bool       `json:"is_new" redis:"is_new"`
	NewFrom          *time.Time `json:"new_from" redis:"new_from"`
	NewTo            *time.Time `json:"new_to" redis:"new_to"`
	Enabled          bool       `json:"enabled" redis:"enabled"`
	ShortDescription string     `json:"short_description" redis:"short_description"`
	FullDescription  string     `json:"full_description" redis:"full_description"`
	Quantity         uint32     `json:"quantity" redis:"quantity"`
	Slug             string     `json:"slug" redis:"slug"`
	CountryOfManf    uint32     `json:"country_of_manufacture" redis:"cty_manf"`
	Weight           *float64   `json:"weight" redis:"weight"`
	BrandId          uint32     `json:"brand_id" redis:"brand_id"`
	ColorId          uint32     `json:"color_id" redis:"color_id"`
	Medias           []uint32   `json:"medias" redis:"medias"`
	MetaTitle        string     `json:"meta_title,omitempty" redis:"meta_title"`
	MetaDescription  string     `json:"meta_description,omitempty" redis:"meta_description"`
	SearchKeywords   *string    `json:"search_keywords" redis:"search_keywords"`
	RelatedProducts  []uint32   `json:"related_products" redis:"related_products"`
	UpSellProducts   []uint32   `json:"upsell" redis:"upsell"`
	CrossSell        []uint32   `json:"cross_sell" redis:"cross_sell"`
	PresetID         *uint32    `json:"preset_id" redis:"preset_id"`
	/** List of category_id */
	PresetValues []uint32 `json:"preset_values" redis:"preset_values"`
}
