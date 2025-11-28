package models

import (
	"database/sql/driver"
	"errors"

	"github.com/goccy/go-json"
)

type AppSettings struct {
	DefaultCurrency        string   `json:"default_currency"`
	SupportedCurrencies    []string `json:"supported_currencies"`
	DefaultLanguage        string   `json:"default_language"`
	SupportedLanguages     []string `json:"supported_languages"`
	SEOMetaTitle           string   `json:"seo_meta_title"`
	SEOMetaDescription     string   `json:"seo_meta_description"`
	EnableSSL              bool     `json:"enable_ssl"`
	DefaultShippingCountry string   `json:"default_shipping_country"`
	EnableRealtimeRates    bool     `json:"enable_realtime_rates"`
	EnabledPaymentGateways []string `json:"enabled_payment_gateways"`
	AllowGuestCheckout     bool     `json:"allow_guest_checkout"`
	TrackInventory         bool     `json:"track_inventory"`
	LowStockThreshold      int      `json:"low_stock_threshold"`
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
