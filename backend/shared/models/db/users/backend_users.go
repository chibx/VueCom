package users

import (
	"time"
)

type Continent struct {
	ID        uint      `gorm:"primarykey" redis:"id"`
	Name      string    `gorm:"not null;unique" redis:"name"`
	Countries []Country `gorm:"foreignKey:ContinentId"`
}

// Country represents a country in the system
type Country struct {
	ID          uint       `gorm:"primarykey" redis:"id"`
	Name        string     `gorm:"not null;unique;index" redis:"name"`
	Code        string     `gorm:"not null;unique;index;type:varchar(5)" redis:"code"`
	ContinentId uint       `gorm:"" redis:"continent_id"`
	Phone       string     `redis:"phone"`
	Currency    string     `redis:"currency"`
	States      []State    `gorm:"foreignKey:CountryID;" redis:"-"`
	Continent   *Continent `gorm:"foreignKey:ContinentId" redis:"-"`
}

type State struct {
	ID        uint   `gorm:"primarykey" redis:"id"`
	Name      string `gorm:"not null;unique;index" redis:"name"`
	CountryID uint   `gorm:"index" redis:"country_id"`
	Cities    []City `gorm:"foreignKey:StateID" redis:"-"`
}

type City struct {
	ID      uint   `gorm:"primarykey" redis:"id"`
	Name    string `gorm:"not null;unique;index" redis:"name"`
	StateID uint   `gorm:"index" redis:"country_id"`
}

func (c City) TableName() string {
	return "cities"
}

type ApiKey struct {
	ID        uint      `gorm:"primarykey" redis:"id"`
	UserID    uint      `gorm:"not null;index" redis:"user_id"`
	KeyPrefix string    `gorm:"type:varchar(16);unique;not null;index" redis:"-"` // Public, indexed identifier
	KeyHash   []byte    `gorm:"type:bytea;not null" redis:"-"`                    // Cryptographic hash of the secret key
	CreatedAt time.Time `gorm:"" redis:"created_at"`
}

type BackendOTP struct {
	UserId     uint         `gorm:"not null"`
	Code       string       `gorm:"not null"`
	ExpiryDate time.Time    `gorm:"not null"`
	User       *BackendUser `gorm:"foreignKey:UserId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" redis:"-"`
}

type BackendRole struct {
	ID           uint     `gorm:"primarykey" redis:"id"`
	Name         string   `gorm:"" redis:"name"`
	ParentID     *uint    `gorm:"" redis:"parent_id"`
	AllowedPerms []string `gorm:"column:allowed_permissions;type:text[]" redis:"allowed_permissions"`
	// ParentRole *BackendRole `gorm:"foreignKey:ParentID" redis:"-"`
}

type BackendUser struct {
	ID              uint                          `gorm:"primarykey" redis:"id"`
	CreatedAt       time.Time                     `gorm:"" redis:"created_at"`
	UpdatedAt       time.Time                     `gorm:"" redis:"updated_at"`
	UserName        string                        `gorm:"column:username;type:varchar(255);index" validate:"" redis:"user_name"`
	FullName        string                        `gorm:"not null;type:varchar(255);index" validate:"" redis:"full_name"`
	Email           string                        `gorm:"unique;not null;type:varchar(255);index" redis:"email"`
	PhoneNumber     *string                       `gorm:"type:varchar(20)" validate:"" redis:"phone_number"`
	Image           *string                       `gorm:"column:image_url" redis:"image_url"`
	CountryId       *uint                         `gorm:"index" redis:"country"`
	Is2FAEnabled    bool                          `gorm:"column:is_2fa_enabled;default:FALSE;not null" redis:"is_2fa_enabled"`
	IsEmailVerified bool                          `gorm:"default:FALSE;not null" redis:"is_email_verified"`
	RoleID          uint                          `gorm:"" redis:"role_id"`
	PasswordHash    string                        `gorm:"not null" redis:"-"`
	ExcludedPerms   []string                      `gorm:"column:excluded_permissions;type:text[]" redis:"excluded_permissions"`
	AdditionalPerms []string                      `gorm:"column:additional_permissions;type:text[]" redis:"additional_permissions"`
	CreatedBy       *uint                         `gorm:"index" redis:"created_by"`
	ByApiKey        bool                          `gorm:"-:all" redis:"by_api_key"` // Track if the user is acting through an API key
	AccountsCreated []BackendUser                 `gorm:"foreignKey:CreatedBy;" redis:"-"`
	Activity        []BackendUserActivity         `gorm:"foreignKey:UserId;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" redis:"-"`
	Sessions        []BackendSession              `gorm:"foreignKey:UserId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" redis:"-"`
	PassResetReqs   []BackendPasswordResetRequest `gorm:"foreignKey:UserId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" redis:"-"`
	Country         *Country                      `gorm:"foreignKey:CountryId;" redis:"-"`
}

type SignupToken struct {
	ID         uint         `gorm:"primarykey"`
	Token      string       `gorm:"not null"`
	Code       string       `gorm:""`
	Supervisor uint         `gorm:""`
	CreatedAt  time.Time    `gorm:"" redis:"created_at"`
	ExpiryAt   time.Time    `gorm:"" redis:"expiry_at"`
	Super      *BackendUser `gorm:"foreignKey:Supervisor"`
}

func (t SignupToken) TableName() string {
	return "backend_signup_tokens"
}

type Backend2FAToken struct {
	UserId            uint         `gorm:"not null;unique" redis:"-"`
	EncryptedToken    string       `gorm:"not null" redis:"-"`
	HashedBackupToken string       `gorm:"not null" redis:"-"`
	User              *BackendUser `gorm:"foreignKey:UserId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" redis:"-"`
}

func (Backend2FAToken) TableName() string {
	return "backend_2fa_tokens"
}

// type BackendSession_Old struct {
// 	UserId    uint         `gorm:"not null" redis:"user_id"`
// 	Token     string       `gorm:"not null" redis:"-"` // redis key would be the token
// 	CreatedAt time.Time    `gorm:"not null" redis:"created_at"`
// 	IpAddr    string       `gorm:"column:ip_address" redis:"ip_address"`
// 	UserAgent string       `gorm:"not null" redis:"user_agent"`
// 	User      *BackendUser `gorm:"foreignKey:UserId;constraint:OnDelete:CASCADE;" redis:"-"`
// }

type BackendSession struct {
	UserId           uint         `gorm:"not null" redis:"user_id"`
	RefreshTokenHash string       `gorm:"not null" redis:"-"`
	LastIP           string       `gorm:"" redis:"last_ip"`
	DeviceId         string       `gorm:"" redis:"device_id"`
	UserAgent        string       `gorm:"not null" redis:"user_agent"`
	CreatedAt        time.Time    `gorm:"not null" redis:"created_at"`
	ExpiresAt        time.Time    `gorm:"not null" redis:"created_at"`
	User             *BackendUser `gorm:"foreignKey:UserId;constraint:OnDelete:CASCADE;" redis:"-"`
}

type BackendUserActivity struct {
	UserId uint `gorm:"index;not null" redis:"user_id"`
	// This would be like "Login", "Password Change", "Profile Update", "Order Handling"
	LogTitle  string    `gorm:"not null" redis:"log_title"`
	Activity  string    `gorm:"not null" redis:"activity"`
	CreatedAt time.Time `gorm:"" redis:"created_at"`
}

func (w BackendUserActivity) TableName() string {
	return "backend_user_activities"
}

type BackendPasswordResetRequest struct {
	ID          uint      `gorm:"primarykey" redis:"id"`
	UserId      uint      `gorm:"not null;index" redis:"user_id"`
	ResetToken  string    `gorm:"not null;unique" redis:"-"` // the key would be the token
	RequestedAt time.Time `gorm:"not null" redis:"requested_at"`
	ExpiresAt   time.Time `gorm:"not null" redis:"expires_at"`
	Used        bool      `gorm:"default:FALSE;not null" redis:"used"`
}

func (w BackendPasswordResetRequest) TableName() string {
	return "password_reset_requests"
}
