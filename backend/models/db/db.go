package db

import (
	"time"
	"vuecom/models"

	"gorm.io/gorm"
)

type ApplicationConfig struct {
	Name       string `gorm:"not null"`
	AdminRoute string `gorm:"not null"`
	Plan       int    `gorm:"not null"`
	ImageUrl   string `gorm:"not null"`
}

type User struct {
	gorm.Model
	models.ApiUser
	PasswordHash string `json:"-" gorm:"not null"`
}

type Customer struct {
	gorm.Model
	FullName        string  `json:"fullName" gorm:"not null;type:varchar(255);index"`
	UserName        *string `json:"userName" gorm:"type:varchar(255);index"`
	PasswordHash    string  `json:"-" gorm:"not null"`
	Email           string  `json:"email" gorm:"unique;not null;type:varchar(255);index"`
	IsEmailVerified bool    `json:"emailVerified" gorm:"default:FALSE;not null"`
	PhoneNumber     *string `json:"phoneNumber" gorm:"type:varchar(20)"`
	Image           *string `json:"image"`
	Country         int     `json:"country" gorm:"index"`
}

type Backend_Session struct {
	gorm.Model
	Token     string    `gorm:"not null"`
	ExpiredAt time.Time `gorm:"not null"`
	IpAddr    string    `gorm:"column:ip_address"`
	UserAgent string    `gorm:"not null"`
	UserID    int       `gorm:"not null"`
}

type Session struct {
	gorm.Model
	Token     string    `gorm:"not null"`
	ExpiredAt time.Time `gorm:"not null"`
	IpAddr    string    `gorm:"column:ip_address"`
	UserAgent string    `gorm:"not null"`
	UserID    int       `gorm:"not null"`
}

// For timed one time password
type OTP struct {
	Code       string    `gorm:"not null"`
	ExpiryDate time.Time `gorm:"not null"`
	UserId     int       `gorm:"not null"`
}

type Product struct {
	gorm.Model
	Name string `gorm:"not null;index;type:text"`
	SKU  string `gorm:"not null;index"`
}

// ID: 131 | Name: Size | Value: XL
// ID: 138 | Name: Color | Value: Black
type Category struct {
	gorm.Model
	Name  string `json:"name" gorm:"type:varchar(255);index"`
	Value string `json:"value" gorm:"type:varchar(100)"`
}
