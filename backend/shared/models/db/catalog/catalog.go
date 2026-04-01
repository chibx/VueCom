package catalog

import (
	"time"
)

// CREATE TYPE promo_code_type AS ENUM ('percentage', 'fixed_amount', 'free_shipping')
type PromoCodeType string
type MediaType string

const (
	PromoCodePercent  PromoCodeType = "percentage"
	PromoCodeFixed    PromoCodeType = "fixed_amount"
	PromoCodeShipping PromoCodeType = "free_shipping"
)

const (
	TYPE_VIDEO MediaType = "vd"
	TYPE_IMAGE MediaType = "img"
)

type Attribute struct {
	ID        uint      `gorm:"primarykey" redis:"id"`
	CreatedAt time.Time `redis:"created_at"`
	UpdatedAt time.Time `redis:"updated_at"`
	Name      string    `gorm:"type:varchar(50);index;not null;unique" redis:"name"`
}

// ID: 131 | Name: Size | Value: XL
// ID: 138 | Name: Color | Value: Black
// TODO: Add a unique constraint on (attribute_id, value)
type Category struct {
	ID          uint       `gorm:"primarykey" redis:"id"`
	CreatedAt   time.Time  `redis:"created_at"`
	UpdatedAt   time.Time  `redis:"updated_at"`
	AttributeID uint       `json:"attributeId" gorm:"index;not null"`
	Value       string     `json:"value" gorm:"index;type:varchar(50);not null"`
	Attribute   *Attribute `json:"-" gorm:"foreignKey:AttributeID"`
}

/**
 * A set of grouped categories that a product should make use of i.e Electronics, Computing
 */
type Preset struct {
	ID         uint               `gorm:"primarykey;index" redis:"id"`
	Name       string             `gorm:"index;type:varchar(50);not null;unique" redis:"name"`
	CreatedAt  time.Time          `gorm:"" redis:"created_at"`
	UpdatedAt  time.Time          `gorm:"" redis:"updated_at"`
	Attributes []PresetAttributes `json:"-" gorm:"foreignKey:PresetID" redis:"-"`
	Products   []Product          `json:"-" gorm:"foreignKey:PresetID" redis:"-"`
}

/**
 * This table would be for linking presets to their respective categories
 * Computing preset could have Color: Red | Storage: 256GB
 * Electornic preset could have Color: Red | Cable Length: 2m
 *
 * TODO: Primary Key (PresetID, CategoryID)
 */
type PresetAttributes struct {
	PresetID   uint      `gorm:"primaryKey;index"`
	CategoryID uint      `gorm:"primaryKey;index"`
	Preset     *Preset   `json:"-" gorm:"foreignKey:PresetID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" redis:"-"`
	Category   *Category `json:"-" gorm:"foreignKey:CategoryID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" redis:"-"`
}

type Tag struct {
	ID   uint   `gorm:"primarykey" redis:"id"`
	Name string `gorm:"not null;unique" redis:"name"`
}

// TODO: Primary Key (ProductID, TagID)
type ProductTags struct {
	ProductID uint `gorm:"primaryKey;autoIncrement:false;"`
	TagID     uint `gorm:"primaryKey;autoIncrement:false;"`
}

type MediaFolder struct {
	ID        uint `gorm:"primaryKey"`
	ParentId  uint
	Name      string
	CreatedAt time.Time
}

type Media struct {
	ID           uint
	FolderId     uint
	Type         MediaType
	Url          string
	ThumbnailUrl string
	CreatedAt    time.Time
	ExternalId   string
	FileName     string
	MimeType     string
	SizeBytes    uint
}

type Product struct {
	ID               uint       `gorm:"primarykey" redis:"id"`
	UpdatedAt        time.Time  `gorm:"" redis:"updated_at"`
	CreatedAt        time.Time  `gorm:"" redis:"created_at"`
	Name             string     `gorm:"not null;index;type:text" redis:"name"`
	SKU              string     `gorm:"not null;index" redis:"sku"`
	BasePrice        float64    `redis:"base_price" gorm:"not null;type:numeric(15, 2)"`
	SalePrice        float64    `redis:"sale_price" gorm:"not null;type:numeric(15, 2)"`
	DiscountStart    *time.Time `redis:"discount_start"`
	DiscountEnd      *time.Time `redis:"discount_end"`
	IsNew            bool       `redis:"is_new"`
	NewFrom          *time.Time `redis:"new_from"`
	NewTo            *time.Time `reids:"new_to"`
	CountryOfManf    uint       `gorm:"column:country_of_manufacture;" redis:"coun_of_manf"`
	Enabled          bool       `redis:"enabled" gorm:"default:TRUE;not null"`
	ShortDescription string     `redis:"short_description"`
	FullDescription  string     `redis:"full_description"`
	Slug             string     `redis:"slug"`
	Weight           *float64   `redis:"weight"`
	BrandId          int        `redis:"brand_id"`
	ColorId          int        `redis:"color_id"`
	MetaTitle        *string    `redis:"meta_title"`
	MetaDescription  *string    `redis:"meta_title"`
	SearchKeywords   *string    `gorm:"column:search_keywords;" redis:"search_keywords"`
	ParentID         *uint      `redis:"parent_id"`
	PresetID         *uint      `gorm:"index" redis:"preset_id"`
	Parent           *Product   `gorm:"foreignKey:ParentID"`
	Preset           *Preset    `gorm:"foreignKey:PresetID;constraint:OnUpdate:SET NULL,OnDelete:SET NULL;" redis:"-"`
	Categories       []Category `gorm:"many2many:product_category_values;" redis:"-"`
	Tags             []Tag      `gorm:"many2many:product_tags;" redis:"-"`
	// Categories  []Category `gorm:"many2many:product_category_values;foreignkey:ID;joinforeignKey:ProductId;References:ID;joinReferences:CategoryId;"`
}

type ProductRelation struct {
	SourceProductID uint64 `gorm:"not null"`
	TargetProductID uint64 `gorm:"not null"`
	RelationType    string `gorm:"not null"`
	SortOrder       int    `gorm:"default:0"`
	CreatedAt       time.Time
}

type ProductMedia struct {
	ProductId uint
	MediaId   uint
	SortOrder uint
	IsMain    bool
}

type ProductCategoryValues struct {
	ProductID  uint `gorm:"primaryKey;autoIncrement:false;"`
	CategoryID uint `gorm:"primaryKey;autoIncrement:false;"`
}

type PromoCode struct {
	ID                 uint             `json:"id" gorm:"primarykey" redis:"id"`
	Name               string           `json:"name" gorm:"not null;index;type:text" redis:"name"`
	Code               string           `json:"code" gorm:"not null;index;unique" redis:"code"`
	Type               PromoCodeType    `json:"type" gorm:"not null" redis:"type"`
	Discount           float64          `json:"discount" gorm:"not null" redis:"discount"`
	MinCartValue       float64          `json:"min_cart_value" gorm:"not null" redis:"min_cart_value"`
	StartDate          time.Time        `json:"start_date" gorm:"" redis:"start_date"`
	ExpiryDate         *time.Time       `json:"expiry_date,omitempty" gorm:"" redis:"expiry_date"`
	UsageLimit         int              `json:"usage_limit" gorm:"not null" redis:"usage_limit"`
	UsageLimitPerUser  int              `json:"usage_limit_per_user" gorm:"not null" redis:"usage_limit_per_user"`
	ProductIDs         []uint           `json:"product_ids" gorm:"type:integer[]" redis:"product_ids"`
	CategoryIDs        []uint           `json:"category_ids" gorm:"type:integer[]" redis:"category_ids"`
	ExcludeProductIDs  []uint           `json:"exclude_product_ids" gorm:"type:integer[]" redis:"exclude_product_ids"`
	ExcludeCategoryIDs []uint           `json:"exclude_category_ids" gorm:"type:integer[]" redis:"exclude_category_ids"`
	IsActive           bool             `json:"is_active" gorm:"default:TRUE;not null" redis:"is_active"`
	UpdatedAt          time.Time        `json:"updated_at" gorm:"" redis:"updated_at"`
	CreatedAt          time.Time        `json:"created_at" gorm:"" redis:"created_at"`
	Usages             []PromoCodeUsage `json:"-" gorm:"foreignKey:CodeID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type PromoCodeUsage struct {
	ID         uint       `gorm:"primarykey" redis:"id"`
	CodeID     uint       `json:"code_id" gorm:"index;not null" redis:"code_id"`
	CustomerID uint       `json:"customer_id" gorm:"index;not null" redis:"customer_id"`
	OrderID    uint       `json:"order_id" gorm:"index;not null" redis:"order_id"`
	UsedAt     time.Time  `json:"used_at" gorm:"index;not null" redis:"used_at"`
	PromoCode  *PromoCode `gorm:"foreignKey:CodeID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" redis:"-"`
}
