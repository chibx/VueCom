package users

import (
	"time"

	"github.com/google/uuid"
)

// For timed one time password
type CustomerOTP struct {
	Code       string    `gorm:"not null"`
	ExpiryDate time.Time `gorm:"not null"`
	CustomerID uint      `gorm:"column:customer_id;not null;index"`
	// Customer   *Customer `gorm:"foreignKey:CustomerID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type CustomerAddress struct {
	ID            uint     `gorm:"primarykey" redis:"id"`
	CustomerID    uint     `gorm:"not null;index" redis:"user_id"`
	StreetAddress string   `gorm:"index;not null" redis:"address"`
	ZipCode       *string  `gorm:"" redis:"zip_code"`
	StateID       *uint    `gorm:"" redis:"state"`
	CountryId     *uint    `gorm:"" redis:"country"`
	State         *State   `gorm:"foreignKey:StateID;" redis:"-"`
	Country       *Country `gorm:"foreignKey:CountryId;" redis:"-"`
	// City          string   `gorm:"index;not null" redis:"city"`
}

type Customer struct {
	ID              uint              `gorm:"primarykey" redis:"id"`
	CreatedAt       time.Time         `gorm:"" redis:"created_at"`
	FullName        string            `gorm:"not null;type:varchar(255);index" redis:"full_name"`
	PasswordHash    *string           `gorm:"" redis:"-"`
	Email           string            `gorm:"unique;not null;type:varchar(255);index" redis:"email"`
	IsEmailVerified bool              `gorm:"default:FALSE;not null" redis:"is_email_verified"`
	PhoneNumber     *string           `gorm:"type:varchar(20);" redis:"phone_number"`
	ImageUrl        *string           `gorm:"column:image_url" redis:"image_url"`
	CountryID       *uint             ``
	BillingAddress  *CustomerAddress  `gorm:"foreignKey:CustomerID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" redis:"-"`
	ShippingAddress *CustomerAddress  `gorm:"foreignKey:CustomerID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" redis:"-"`
	Sessions        []CustomerSession `gorm:"foreignKey:CustomerID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" redis:"-"`
	OTPs            []CustomerOTP     `gorm:"foreignKey:CustomerID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" redis:"-"`
	Wishlist        []WishlistItem    `gorm:"foreignKey:CustomerID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" redis:"-"`
	Cart            []CartItem        `gorm:"foreignKey:CustomerID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" redis:"-"`
	// Wishlist        []Product         `gorm:"many2many:customer.customer_wishlists;"`
}

// type CustomerSession_Old struct {
// 	UserID    uint      `gorm:"not null" redis:"user_id"`
// 	Token     string    `gorm:"not null" redis:"-"` // redis key would be the token
// 	ExpiredAt time.Time `gorm:"not null" redis:"expired_at"`
// 	IpAddr    string    `gorm:"column:ip_address" redis:"ip_address"`
// 	UserAgent string    `gorm:"not null" redis:"user_agent"`
// 	Customer  *Customer `gorm:"foreignKey:CustomerID;constraint:OnDelete:CASCADE;" redis:"-"`
// }

type CustomerSession struct {
	ID               uuid.UUID    `gorm:"primarykey;not null;type:uuid" redis:"-"` // redis key would be the token id
	UserId           uint         `gorm:"not null" redis:"user_id"`
	RefreshTokenHash string       `gorm:"not null" redis:"-"`
	LastIP           string       `gorm:"" redis:"last_ip"`
	DeviceId         string       `gorm:"" redis:"device_id"`
	UserAgent        string       `gorm:"not null" redis:"user_agent"`
	CreatedAt        time.Time    `gorm:"not null" redis:"created_at"`
	ExpiresAt        time.Time    `gorm:"not null" redis:"created_at"`
	User             *BackendUser `gorm:"foreignKey:UserId;constraint:OnDelete:CASCADE;" redis:"-"`
}

type CartItem struct {
	CustomerID uint      `gorm:"index"`
	ProductID  uint      `gorm:"index"`
	Quantity   int       `gorm:"not null;default:1"`
	AddedAt    time.Time `gorm:""`
	// Product    *Product  `gorm:"foreignKey:ProductID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

func (c CartItem) TableName() string {
	return "customer_carts"
}

type WishlistItem struct {
	CustomerID uint      `gorm:"index"`
	ProductID  uint      `gorm:"index"`
	AddedAt    time.Time `gorm:"not null"`
	// Product    *Product  `gorm:"foreignKey:ProductID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
