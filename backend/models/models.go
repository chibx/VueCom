package models

import (
	"time"

	"gorm.io/gorm"
)

type BackendUser struct {
	gorm.Model
	FullName        string
	UserName        *string
	PasswordHash    string
	Email           string `gorm:"unique;"`
	IsEmailVerified bool   `gorm:"default:FALSE"`
	PhoneNumber     *string
	DateOfBirth     string
	Image           string
	Country         int
	Plan            int
}

type Session struct {
	gorm.Model
	Token     string
	ExpiredAt time.Time
	IpAddr    string `gorm:"column:ip_address"`
	UserAgent string
	UserId    int
}

// For timed one time password
type OTP struct {
	Code       string
	ExpiryDate time.Time
	UserId     int
}

// JWT Format sent back to the client
type BackendJWTPayload struct {
	UserId int    `json:"userId"`
	Role   string `json:"role"`
	Exp    int    `json:"exp"`
}
