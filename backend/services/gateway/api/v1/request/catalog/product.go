package catalog

import "time"

type ProductVisibility int

const (
	VisibilityCatalog ProductVisibility = 0
	VisibilitySearch  ProductVisibility = 1
	VisibilityBoth    ProductVisibility = 2
)

type CreateProductReq struct {
	Name             string     `json:"name" validate:"required" name:"Product Name"`
	SKU              string     `json:"sku" validate:"required" name:"Product SKU"`
	BasePrice        float64    `json:"base_price" validate:"required" name:"Base Price"`
	SalePrice        float64    `json:"sale_price" validate:"required" name:"Discount Price"`
	InStock          bool       `json:"in_stock"`
	DiscountStart    *time.Time `json:"discount_start"`
	DiscountEnd      *time.Time `json:"discount_end"`
	Categories       []int32    `json:"categories"`
	IsNew            bool       `json:"is_new"`
	NewFrom          *time.Time `json:"new_from"`
	NewTo            *time.Time `json:"new_to"`
	Enabled          bool       `json:"enabled"`
	ShortDescription string     `json:"short_description"`
	FullDescription  string     `json:"full_description"`
	Quantity         uint32     `json:"quantity"`
	Slug             string     `json:"slug" validate:"required" name:"Product Slug"`
	CountryOfManf    uint32     `json:"country_of_manufacture" validate:"required" name:"Country of Manufacture"`
	Weight           *float64   `json:"weight"`
	BrandId          uint32     `json:"brand_id"`
	ColorId          uint32     `json:"color_id"`
	Medias           []uint32   `json:"medias"`
	MetaTitle        string     `json:"meta_title,omitempty" validate:"required" name:"Meta Title"`
	MetaDescription  string     `json:"meta_description,omitempty" validate:"required,maxlength=150" name:"Meta Description"`
	SearchKeywords   *string    `json:"search_keywords"`
	RelatedProducts  []uint32   `json:"related_products"`
	UpSellProducts   []uint32   `json:"upsell"`
	CrossSell        []uint32   `json:"cross_sell"`
	PresetID         *uint32    `json:"preset_id" validate:"omitnil,gte=0" name:"Product Preset"`
	/** List of category_id */
	PresetValues []uint32 `json:"preset_values"`
}
