package product

import "time"

type ApiProducts struct {
	Name string `json:"name" gorm:"not null;index;type:text"`
	SKU  string `json:"sku" gorm:"not null;index"`
	// Just made the precision to be 15 (Don't know how bad the ecomomy of some countries are)
	Price      float64   `json:"price" gorm:"not null;type:numeric(15, 2)"`
	DscPercent float64   `json:"dsc_percent" gorm:"type:numeric(5, 2)"`
	DscPeriod  time.Time `json:"dsc_period" gorm:""`
	Enabled    bool      `json:"enabled" gorm:""`
	// Warranty   string    `json:"warranty"`
	Description string
	Url         string `json:"url"`
}
