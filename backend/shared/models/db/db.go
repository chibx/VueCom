package db

import (
	"time"

	"vuecom/shared/models"
)

type ApplicationConfig struct {
	Name string `gorm:"not null"`
	// Dashboard
	AdminRoute string `gorm:"not null"`
	Plan       uint   `gorm:"not null"`
	// Logo
	ImageUrl string `gorm:"not null"`
}

type BackendUser struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	models.ApiBackendUser
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
	models.ApiProducts
}

type ApiKey struct {
	ID           uint      `gorm:"primarykey" redis:"id"`
	UserID       uint      `gorm:"not null;index" redis:"user_id"`
	EncryptedKey []byte    `gorm:"type:bytea;not null" redis:"encrypted_key"`
	CreatedAt    time.Time `gorm:"" redis:"created_at"`
}
