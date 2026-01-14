package users

import (
	"time"
)

// Country represents a country in the system
type Country struct {
	ID        uint      `gorm:"primarykey" redis:"id"`
	Name      string    `gorm:"not null;unique;index" redis:"name"`
	Code      string    `gorm:"not null;unique;index;type:varchar(5)" redis:"code"`
	CreatedAt time.Time `gorm:"" redis:"created_at"`
	States    []State   `gorm:"foreignKey:CountryID;"`
}

func (Country) TableName() string {
	return "backend.countries"
}

type State struct {
	ID        uint   `gorm:"primarykey" redis:"id"`
	Name      string `gorm:"not null;unique;index" redis:"name"`
	CountryID uint   `gorm:"index" redis:"country_id"`
}

func (State) TableName() string {
	return "backend.states"
}

// type AppData struct {
// 	Name       string             `json:"app_name" gorm:"" redis:"name"`
// 	AdminRoute string             `json:"-" gorm:"" redis:"admin_route"`
// 	LogoUrl    string             `json:"app_logo" gorm:"" redis:"logo_url"`
// 	Settings   models.AppSettings `json:"settings" gorm:"type:jsonb;default:\"{}\"" redis:"settings"`
// }

// func (AppData) TableName() string {
// 	return "backend.app_data"
// }

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
	ID              uint                          `gorm:"primarykey" redis:"id"`
	CreatedAt       time.Time                     `gorm:"" redis:"created_at"`
	UpdatedAt       time.Time                     `gorm:"" redis:"updated_at"`
	UserName        *string                       `gorm:"type:varchar(255);index" validate:"" redis:"user_name"`
	FullName        string                        `gorm:"not null;type:varchar(255);index" validate:"" redis:"full_name"`
	Email           string                        `gorm:"unique;not null;type:varchar(255);index" redis:"email"`
	PhoneNumber     *string                       `gorm:"type:varchar(20)" validate:"" redis:"phone_number"`
	Image           *string                       `gorm:"column:image_url" redis:"image_url"`
	CountryId       *uint                         `gorm:"index" redis:"country"`
	IsEmailVerified bool                          `gorm:"default:FALSE;not null" redis:"-"`
	Role            string                        `gorm:"type:varchar(50)" redis:"role"`
	PasswordHash    string                        `gorm:"not null" redis:"-"`
	CreatedBy       *uint                         `gorm:"index" redis:"created_by"`
	AccountsCreated []BackendUser                 `gorm:"foreignKey:CreatedBy;" redis:"-"`
	ByApiKey        bool                          `gorm:"-:all" redis:"by_api_key"` // Track if the user is acting through an API key
	Activity        []BackendUserActivity         `gorm:"foreignKey:UserId;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" redis:"-"`
	Sessions        []BackendSession              `gorm:"foreignKey:UserId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" redis:"-"`
	OTP             []BackendOTP                  `gorm:"foreignKey:UserId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" redis:"-"`
	PassResetReqs   []BackendPasswordResetRequest `gorm:"foreignKey:UserId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" redis:"-"`
	Country         *Country                      `gorm:"foreignKey:CountryId;" redis:"-"`
}

func (BackendUser) TableName() string {
	return "backend.backend_users"
}

type BackendSession struct {
	UserId    uint         `gorm:"not null" redis:"user_id"`
	Token     string       `gorm:"not null" redis:"-"` // redis key would be the token
	CreatedAt time.Time    `gorm:"not null" redis:"created_at"`
	IpAddr    string       `gorm:"column:ip_address" redis:"ip_address"`
	UserAgent string       `gorm:"not null" redis:"user_agent"`
	User      *BackendUser `gorm:"foreignKey:UserId;constraint:OnDelete:CASCADE;" redis:"-"`
}

func (BackendSession) TableName() string {
	return "backend.backend_sessions"
}

type BackendUserActivity struct {
	UserId uint `gorm:"index;not null" redis:"user_id"`
	// This would be like "Login", "Password Change", "Profile Update", "Order Handling"
	LogTitle  string    `gorm:"type:varchar(100);not null" redis:"log_title"`
	Activity  string    `gorm:"type:text;not null" redis:"activity"`
	CreatedAt time.Time `gorm:"" redis:"created_at"`
}

func (BackendUserActivity) TableName() string {
	return "backend.backend_user_activities"
}

type BackendPasswordResetRequest struct {
	Id          uint      `gorm:"primarykey" redis:"id"`
	UserId      uint      `gorm:"not null;index" redis:"user_id"`
	ResetToken  string    `gorm:"not null;unique" redis:"-"` // the key would be the token
	RequestedAt time.Time `gorm:"not null" redis:"requested_at"`
	ExpiresAt   time.Time `gorm:"not null" redis:"expires_at"`
	Used        bool      `gorm:"default:FALSE;not null" redis:"used"`
}

func (BackendPasswordResetRequest) TableName() string {
	return "backend.backend_password_reset_requests"
}
