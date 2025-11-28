package db

import (
	"time"
	"vuecom/shared/models"
)

type AppData struct {
	Name       string             `json:"app_name" gorm:"" redis:"name"`
	AdminRoute string             `json:"-" gorm:"" redis:"admin_route"`
	LogoUrl    string             `json:"app_logo" gorm:"" redis:"logo_url"`
	Settings   models.AppSettings `json:"settings" gorm:"type:jsonb;default:\"{}\"" redis:"settings"`
}

func (AppData) TableName() string {
	return "backend.app_data"
}

type ApiKey struct {
	ID        uint      `gorm:"primarykey" redis:"id"`
	UserID    uint      `gorm:"not null;index" redis:"user_id"`
	KeyPrefix string    `gorm:"type:varchar(16);unique;not null;index" redis:"-"` // Public, indexed identifier
	KeyHash   []byte    `gorm:"type:bytea;not null" redis:"-"`                    // Cryptographic hash of the secret key
	CreatedAt time.Time `gorm:"" redis:"created_at"`
}

func (ApiKey) TableName() string {
	return "backend.api_keys"
}

type BackendOTP struct {
	Code       string    `gorm:"not null"`
	ExpiryDate time.Time `gorm:"not null"`
	UserId     uint      `gorm:"not null"`
}

func (BackendOTP) TableName() string {
	return "backend.backend_otps"
}

type BackendUser struct {
	ID              uint `gorm:"primarykey"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
	UserName        *string                `gorm:"type:varchar(255);index" validate:""`
	FullName        string                 `gorm:"not null;type:varchar(255);index" validate:""`
	Email           string                 `gorm:"unique;not null;type:varchar(255);index"`
	PhoneNumber     *string                `gorm:"type:varchar(20)" validate:""`
	Image           *string                `gorm:"column:image_url"`
	Country         *uint                  `gorm:"index"`
	IsEmailVerified bool                   `gorm:"default:FALSE;not null"`
	Role            string                 `gorm:"type:varchar(50)"`
	PasswordHash    string                 `gorm:"not null"`
	CreatedBy       uint                   `gorm:"index,not null"`
	Activity        *[]BackendUserActivity `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Sessions        *[]BackendSession      `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	OTP             *[]BackendOTP          `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

func (BackendUser) TableName() string {
	return "backend.backend_users"
}

type BackendSession struct {
	Token     string    `gorm:"not null" redis:"token"`
	ExpiredAt time.Time `gorm:"not null" redis:"expired_at"`
	IpAddr    string    `gorm:"column:ip_address" redis:"ip_address"`
	UserAgent string    `gorm:"not null" redis:"user_agent"`
	UserID    uint      `gorm:"not null" redis:"user_id"`
}

func (BackendSession) TableName() string {
	return "backend.backend_sessions"
}

type BackendUserActivity struct {
	UserID uint `gorm:"index;not null"`
	// This would be like "Login", "Password Change", "Profile Update", "Order Handling"
	LogTitle  string    `gorm:"type:varchar(100);not null"`
	Activity  string    `gorm:"type:text;not null"`
	CreatedAt time.Time `gorm:""`
}

func (BackendUserActivity) TableName() string {
	return "backend.backend_user_activities"
}
