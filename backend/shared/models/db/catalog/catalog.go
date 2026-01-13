package catalog

import (
	"time"
)

type Attribute struct {
	ID        uint      `gorm:"primarykey" redis:"id"`
	CreatedAt time.Time `redis:"created_at"`
	UpdatedAt time.Time `redis:"updated_at"`
	Name      string    `gorm:"type:varchar(50);index;not null;unique" redis:"name"`
}

func (Attribute) TableName() string {
	return "catalog.attributes"
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
	Attribute   *Attribute `gorm:"foreignKey:AttributeID"`
}

func (Category) TableName() string {
	return "catalog.category"
}

/**
 * A set of grouped categories that a product should make use of i.e Electronics, Computing
 */
type Preset struct {
	ID         uint               `gorm:"primarykey;index" redis:"id"`
	Name       string             `gorm:"index;type:varchar(50);not null;unique" redis:"name"`
	CreatedAt  time.Time          `gorm:"" redis:"created_at"`
	UpdatedAt  time.Time          `gorm:"" redis:"updated_at"`
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
	Preset     *Preset   `gorm:"foreignKey:PresetID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" redis:"-"`
	Category   *Category `gorm:"foreignKey:CategoryID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" redis:"-"`
}

func (PresetAttributes) TableName() string {
	return "catalog.preset_attributes"
}

type Tag struct {
	ID   uint   `gorm:"primarykey" redis:"id"`
	Name string `gorm:"not null;unique" redis:"name"`
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
	Name        string    `json:"name" gorm:"not null;index;type:text" redis:"name"`
	SKU         string    `json:"sku" gorm:"not null;index" redis:"sku"`
	Price       float64   `json:"price" gorm:"not null;type:numeric(15, 2)" redis:"price"`
	DscPercent  float64   `json:"dsc_percent" gorm:"type:numeric(5, 2)"`
	DscPeriod   time.Time `json:"dsc_period" gorm:""`
	Enabled     bool      `json:"enabled" gorm:"default:TRUE;not null"`
	Description string    `json:"description"`
	Url         string    `json:"url"`
	ImageUrl    *string   `gorm:"column:image_url"`
	PresetID    *uint     `gorm:"index"`
	Preset      *Preset   `gorm:"foreignKey:PresetID;constraint:OnUpdate:SET NULL,OnDelete:SET NULL;" redis:"-"`
	// Categories  []Category `gorm:"many2many:catalog.product_category_values;foreignkey:ID;joinforeignKey:ProductId;References:ID;joinReferences:CategoryId;"`
	Categories []Category `gorm:"many2many:catalog.product_category_values;" redis:"-"`
	Tags       []Tag      `gorm:"many2many:catalog.product_tags;" redis:"-"`
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
	Name               string           `json:"name" gorm:"not null;index;type:text" redis:"name"`
	Code               string           `json:"code" gorm:"not null;index;unique" redis:"code"`
	Type               string           `json:"type" gorm:"not null" redis:"type"`
	Discount           float64          `json:"discount" gorm:"not null" redis:"discount"`
	MinCartValue       float64          `json:"min_cart_value" gorm:"not null" redis:"min_cart_value"`
	ExpiryDate         time.Time        `json:"expiry_date" gorm:"" redis:"expiry_date"`
	StartDate          time.Time        `json:"start_date" gorm:"" redis:"start_date"`
	UsageLimit         int              `json:"usage_limit" gorm:"not null" redis:"usage_limit"`
	UsageLimitPerUser  int              `json:"usage_limit_per_user" gorm:"not null" redis:"usage_limit_per_user"`
	ProductIDs         []uint           `json:"product_ids" gorm:"type:jsonb" redis:"product_ids"`
	CategoryIDs        []uint           `json:"category_ids" gorm:"type:jsonb" redis:"category_ids"`
	ExcludeProductIDs  []uint           `json:"exclude_product_ids" gorm:"type:jsonb" redis:"exclude_product_ids"`
	ExcludeCategoryIDs []uint           `json:"exclude_category_ids" gorm:"type:jsonb" redis:"exclude_category_ids"`
	IsActive           bool             `json:"is_active" gorm:"default:TRUE;not null" redis:"is_active"`
	Usages             []PromoCodeUsage `gorm:"foreignKey:CodeID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

func (PromoCode) TableName() string {
	return "catalog.promo_codes"
}

type PromoCodeUsage struct {
	ID        uint       `gorm:"primarykey" redis:"id"`
	CodeID    uint       `json:"code_id" gorm:"index;not null" redis:"code_id"`
	UserID    uint       `json:"user_id" gorm:"index;not null" redis:"user_id"`
	OrderID   uint       `json:"order_id" gorm:"index;not null" redis:"order_id"`
	UsedAt    time.Time  `json:"used_at" gorm:"" redis:"used_at"`
	PromoCode *PromoCode `gorm:"foreignKey:CodeID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" redis:"-"`
}

func (PromoCodeUsage) TableName() string {
	return "catalog.promo_code_usages"
}
