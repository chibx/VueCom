package product

import "time"

type ApiProducts struct {
	Name string `json:"name"`
	SKU  string `json:"sku"`
	// Just made the precision to be 15 (Don't know how bad the ecomomy of some countries are)
	Price      float64   `json:"price"`
	DscPercent float64   `json:"dsc_percent"`
	DscPeriod  time.Time `json:"dsc_period"`
	Enabled    bool      `json:"enabled"`
	// Warranty   string    `json:"warranty"`
	Description string
	Url         string `json:"url"`
}
