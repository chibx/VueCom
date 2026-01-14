package appdata

import (
	"database/sql/driver"
	"errors"

	"github.com/goccy/go-json"
)

type AppData struct {
	Name       string      `json:"app_name" gorm:"" redis:"name"`
	AdminRoute string      `json:"-" gorm:"" redis:"admin_route"`
	LogoUrl    string      `json:"app_logo" gorm:"" redis:"logo_url"`
	Settings   AppSettings `json:"settings" gorm:"type:jsonb;default:\"{}\"" redis:"settings"`
}

type AppSettings struct {
	DefaultCurrency        string   `json:"default_currency" validate:"required"`
	SupportedCurrencies    []string `json:"supported_currencies" validate:"required"`
	DefaultLanguage        string   `json:"default_language" validate:"required"`
	SupportedLanguages     []string `json:"supported_languages" validate:"required"`
	SEOMetaTitle           string   `json:"seo_meta_title" validate:""`
	SEOMetaDescription     string   `json:"seo_meta_description" validate:""`
	EnableSSL              bool     `json:"enable_ssl" validate:""`
	DefaultShippingCountry string   `json:"default_shipping_country" validate:"required,min=2,max=5"`
	EnableRealtimeRates    bool     `json:"enable_realtime_rates" validate:""`
	EnabledPaymentGateways []string `json:"enabled_payment_gateways" validate:"required"`
	AllowGuestCheckout     bool     `json:"allow_guest_checkout" validate:""`
	TrackInventory         bool     `json:"track_inventory" validate:""`
	LowStockThreshold      int      `json:"low_stock_threshold" validate:"gt=0"`
}

func (c AppSettings) Value() (driver.Value, error) {
	b, err := json.Marshal(c)
	return string(b[:]), err
}

func (c *AppSettings) Scan(src any) error {
	var value AppSettings
	str, ok := src.(string)
	if !ok {
		byt, ok := src.([]byte)
		if !ok {
			return errors.New("Embedded.Scan byte assertion failed")
		}
		if err := json.Unmarshal(byt, &value); err != nil {
			return err
		}
	} else {
		if err := json.Unmarshal([]byte(str), &value); err != nil {
			return err
		}
	}
	*c = value
	return nil
}

func GetDefaultAppSettings() AppSettings {
	return AppSettings{
		DefaultCurrency:        "USD",
		SupportedCurrencies:    []string{"USD", "EUR", "GBP"},
		DefaultLanguage:        "en",
		SupportedLanguages:     []string{"en", "fr", "es"},
		SEOMetaTitle:           "Vuecom - E-commerce Platform",
		SEOMetaDescription:     "A modern e-commerce platform.",
		EnableSSL:              true,
		DefaultShippingCountry: "US",
		EnableRealtimeRates:    false,
		EnabledPaymentGateways: []string{"stripe", "paypal"},
		AllowGuestCheckout:     true,
		TrackInventory:         true,
		LowStockThreshold:      10,
	}
}
