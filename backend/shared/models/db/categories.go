package db

import (
	"time"
)

type Attribute struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Name      string `gorm:"type:varchar(50);index;not null;unique"`
}

func (Attribute) TableName() string {
	return "catalog.attributes"
}

// ID: 131 | Name: Size | Value: XL
// ID: 138 | Name: Color | Value: Black
// TODO: Add a unique constraint on (attribute_id, value)
type Category struct {
	ID          uint `gorm:"primarykey"`
	CreatedAt   time.Time
	AttributeID uint       `json:"attributeId" gorm:"index;not null"`
	Value       string     `json:"value" gorm:"index;type:varchar(50);not null"`
	Attribute   *Attribute `gorm:"foreignKey:AttributeID"`
}

func (Category) TableName() string {
	return "catalog.category"
}

/**
 * A set of grouped categories that a product should make use of i.e Electronics, Computing
 */
type Preset struct {
	ID         uint               `gorm:"primarykey;index"`
	Name       string             `gorm:"index;type:varchar(50);not null;unique"`
	CreatedAt  time.Time          `gorm:""`
	UpdatedAt  time.Time          `gorm:""`
	Attributes []PresetAttributes `gorm:"foreignKey:PresetID"`
	Products   []Product          `gorm:"foreignKey:PresetID"`
}

func (Preset) TableName() string {
	return "catalog.presets"
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
	Preset     *Preset   `gorm:"foreignKey:PresetID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Category   *Category `gorm:"foreignKey:CategoryID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

func (PresetAttributes) TableName() string {
	return "catalog.preset_attributes"
}

type Tag struct {
	ID   uint   `gorm:"primarykey"`
	Name string `gorm:"not null;unique"`
}

func (Tag) TableName() string {
	return "catalog.tags"
}

// TODO: Primary Key (ProductID, TagID)
type ProductTags struct {
	ProductID uint `gorm:"primaryKey;autoIncrement:false;"`
	TagID     uint `gorm:"primaryKey;autoIncrement:false;"`
}

func (ProductTags) TableName() string {
	return "catalog.product_tags"
}

type Product struct {
	ID          uint      `gorm:"primarykey" redis:"id"`
	UpdatedAt   time.Time `gorm:"" redis:"updated_at"`
	CreatedAt   time.Time `gorm:"" redis:"created_at"`
	Name        string    `json:"name" gorm:"not null;index;type:text"`
	SKU         string    `json:"sku" gorm:"not null;index"`
	Price       float64   `json:"price" gorm:"not null;type:numeric(15, 2)"`
	DscPercent  float64   `json:"dsc_percent" gorm:"type:numeric(5, 2)"`
	DscPeriod   time.Time `json:"dsc_period" gorm:""`
	Enabled     bool      `json:"enabled" gorm:"default:TRUE;not null"`
	Description string    `json:"description"`
	Url         string    `json:"url"`
	ImageUrl    *string   `gorm:"column:image_url"`
	PresetID    *uint     `gorm:"index"`
	Preset      *Preset   `gorm:"foreignKey:PresetID;constraint:OnUpdate:SET NULL,OnDelete:SET NULL;"`
	// Categories  []Category `gorm:"many2many:catalog.product_category_values;foreignkey:ID;joinforeignKey:ProductId;References:ID;joinReferences:CategoryId;"`
	Categories []Category `gorm:"many2many:catalog.product_category_values;"`
	Tags       []Tag      `gorm:"many2many:catalog.product_tags;"`
}

func (Product) TableName() string {
	return "catalog.products"
}

type ProductCategoryValues struct {
	ProductID  uint `gorm:"primaryKey;autoIncrement:false;"`
	CategoryID uint `gorm:"primaryKey;autoIncrement:false;"`
}

func (ProductCategoryValues) TableName() string {
	return "catalog.product_category_values"
}

type PromoCode struct {
	ID                 uint             `gorm:"primarykey" redis:"id"`
	UpdatedAt          time.Time        `gorm:"" redis:"updated_at"`
	CreatedAt          time.Time        `gorm:"" redis:"created_at"`
	Name               string           `json:"name" gorm:"not null;index;type:text"`
	Code               string           `json:"code" gorm:"not null;index;unique"`
	Type               string           `json:"type" gorm:"not null"`
	Discount           float64          `json:"discount" gorm:"not null"`
	MinCartValue       float64          `json:"min_cart_value" gorm:"not null"`
	ExpiryDate         time.Time        `json:"expiry_date" gorm:""`
	StartDate          time.Time        `json:"start_date" gorm:""`
	UsageLimit         int              `json:"usage_limit" gorm:"not null"`
	UsageLimitPerUser  int              `json:"usage_limit_per_user" gorm:"not null"`
	ProductIDs         []uint           `json:"product_ids" gorm:"type:jsonb"`
	CategoryIDs        []uint           `json:"category_ids" gorm:"type:jsonb"`
	ExcludeProductIDs  []uint           `json:"exclude_product_ids" gorm:"type:jsonb"`
	ExcludeCategoryIDs []uint           `json:"exclude_category_ids" gorm:"type:jsonb"`
	IsActive           bool             `json:"is_active" gorm:"default:TRUE;not null"`
	Usages             []PromoCodeUsage `gorm:"foreignKey:CodeID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

func (PromoCode) TableName() string {
	return "catalog.promo_codes"
}

type PromoCodeUsage struct {
	ID        uint       `gorm:"primarykey" redis:"id"`
	CodeID    uint       `json:"code_id" gorm:"index;not null"`
	UserID    uint       `json:"user_id" gorm:"index;not null"`
	OrderID   uint       `json:"order_id" gorm:"index;not null"`
	UsedAt    time.Time  `json:"used_at" gorm:""`
	PromoCode *PromoCode `gorm:"foreignKey:CodeID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Customer  *Customer  `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Order     *Order     `gorm:"foreignKey:OrderID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

func (PromoCodeUsage) TableName() string {
	return "catalog.promo_code_usages"
}
