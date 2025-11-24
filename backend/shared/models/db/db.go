package db

import (
	"time"
)

type AppData struct {
	Name       string `json:"app_name" gorm:"" redis:"name"`
	AdminRoute string `json:"-" gorm:"" redis:"admin_route"`
	LogoUrl    string `json:"app_logo" gorm:"" redis:"logo_url"`
	// Plan       int    `json:"app_plan"`
}

type sharedUserProps struct {
	FullName        string  `json:"full_name" gorm:"not null;type:varchar(255);index" validate:""`
	UserName        *string `json:"user_name" gorm:"type:varchar(255);index" validate:""`
	Email           string  `json:"email" gorm:"unique;not null;type:varchar(255);index"`
	PhoneNumber     *string `json:"phone_number" gorm:"type:varchar(20)" validate:""`
	Image           *string `json:"image" validate:"url"`
	Country         uint    `json:"country" gorm:"index"`
	IsEmailVerified bool    `json:"email_verified" gorm:"default:FALSE;not null"`
	Password        *string `json:"password,omitempty" validate:"required"`
}

type BackendUser struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	sharedUserProps
	Role         string `json:"role" gorm:"type:varchar(50)"`
	PasswordHash string `json:"-" gorm:"not null"`
}

type Customer struct {
	ID              uint `gorm:"primarykey"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
	FullName        string  `json:"fullName" gorm:"not null;type:varchar(255);index"`
	UserName        *string `json:"userName" gorm:"type:varchar(255);index"`
	PasswordHash    *string `json:"-" gorm:""`
	Email           string  `json:"email" gorm:"unique;not null;type:varchar(255);index"`
	IsEmailVerified bool    `json:"emailVerified" gorm:"default:FALSE;not null"`
	PhoneNumber     *string `json:"phoneNumber" gorm:"type:varchar(20)"`
	Image           *string `json:"image"`
	Country         uint    `json:"country" gorm:"index"`
}

type BackendSession struct {
	ID        uint      `gorm:"primarykey"`
	Token     string    `gorm:"not null"`
	ExpiredAt time.Time `gorm:"not null"`
	IpAddr    string    `gorm:"column:ip_address"`
	UserAgent string    `gorm:"not null"`
	UserID    uint      `gorm:"not null"`
}

type CustomerSession struct {
	ID        uint      `gorm:"primarykey" redis:"id"`
	Token     string    `gorm:"not null" redis:"token"`
	ExpiredAt time.Time `gorm:"not null" redis:"expired_at"`
	IpAddr    string    `gorm:"column:ip_address" redis:"ip_addr"`
	UserAgent string    `gorm:"not null" redis:"user_agent"`
	UserID    uint      `gorm:"not null" redis:"user_id"`
}

// For timed one time password
type OTP struct {
	Code       string    `gorm:"not null"`
	ExpiryDate time.Time `gorm:"not null"`
	UserId     uint      `gorm:"not null"`
}

type Product struct {
	ID        uint      `gorm:"primarykey" redis:"id"`
	UpdatedAt time.Time `redis:"updated_at"`
	CreatedAt time.Time `redis:"created_at"`
	Name      string    `json:"name" gorm:"not null;index;type:text"`
	SKU       string    `json:"sku" gorm:"not null;index"`
	// Just made the precision to be 15 (Don't know how bad the ecomomy of some countries are)
	Price      float64   `json:"price" gorm:"not null;type:numeric(15, 2)"`
	DscPercent float64   `json:"dsc_percent" gorm:"type:numeric(5, 2)"`
	DscPeriod  time.Time `json:"dsc_period" gorm:""`
	Enabled    bool      `json:"enabled" gorm:""`
	// Warranty   string    `json:"warranty"`
	Description string
	Url         string `json:"url"`
	// models.ApiProducts
}

type ApiKey struct {
	ID           uint      `gorm:"primarykey" redis:"id"`
	UserID       uint      `gorm:"not null;index" redis:"user_id"`
	EncryptedKey []byte    `gorm:"type:bytea;not null" redis:"encrypted_key"`
	CreatedAt    time.Time `gorm:"" redis:"created_at"`
}
