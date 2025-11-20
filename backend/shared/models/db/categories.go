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

// ID: 131 | Name: Size | Value: XL
// ID: 138 | Name: Color | Value: Black
// TODO: Add a unique constraint on (attribute_id, value)
type Category struct {
	ID          uint `gorm:"primarykey"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	AttributeID uint   `json:"attributeId" gorm:"index;not null"`
	Value       string `json:"value" gorm:"index;type:varchar(50);not null"`
}

/**
 * A set of grouped categories that a product should make use of i.e Electronics, Computing
 */
type Preset struct {
	ID   uint   `gorm:"primarykey;index"`
	Name string `gorm:"index;type:varchar(50);not null;unique"`
}

/**
 * This table would be for linking presets to their respective categories
 * Computing preset could have Color: Red | Storage: 256GB
 * Electornic preset could have Color: Red | Cable Length: 2m
 *
 * TODO: Primary Key (PresetID, CategoryID)
 */
type PresetAttributes struct {
	PresetID   uint `gorm:""`
	CategoryID uint `gorm:""`
}

type Tag struct {
	ID   uint   `gorm:"primarykey"`
	Name string `gorm:"not null;unique"`
}

// TODO: Primary Key (ProductID, TagID)
type ProductTags struct {
	ProductID uint
	TagID     uint
}
