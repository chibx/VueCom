package db

import "time"

// For timed one time password
type CustomerOTP struct {
	Code       string    `gorm:"not null"`
	ExpiryDate time.Time `gorm:"not null"`
	CustomerID uint      `gorm:"column:customer_id;not null;index"`
	// Customer   *Customer `gorm:"foreignKey:CustomerID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

func (CustomerOTP) TableName() string {
	return "customer.customer_otps"
}

type Customer struct {
	ID              uint      `gorm:"primarykey" redis:"id"`
	CreatedAt       time.Time `gorm:"" redis:"created_at"`
	FullName        string    `gorm:"not null;type:varchar(255);index" redis:"full_name"`
	UserName        *string   `gorm:"type:varchar(255);index" redis:"user_name"`
	PasswordHash    *string   `gorm:"" redis:"-"`
	Email           string    `gorm:"unique;not null;type:varchar(255);index" redis:"email"`
	IsEmailVerified bool      `gorm:"default:FALSE;not null" redis:"is_email_verified"`
	PhoneNumber     string    `gorm:"type:varchar(20);not null" redis:"phone_number"`
	Image           *string   `gorm:"column:image_url" redis:"image_url"`
	CountryID       uint      `gorm:"column:country;index;not null" redis:"country"`
	Country         *Country  `gorm:"foreignKey:CountryID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	City            string    `gorm:"index;not null" redis:"city"`
	State           string    `gorm:"index;not null" redis:"state"`
	ZipCode         string    `gorm:"index;not null" redis:"zip_code"`
	Address         string    `gorm:"index;not null" redis:"address"`
	// Sessions        []CustomerSession `gorm:"foreignKey:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	// OTPs            []CustomerOTP     `gorm:"foreignKey:CustomerID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	// Wishlist        []*Product        `gorm:"many2many:customer.customer_wishlists;"`
	// Cart            []*CartItem       `gorm:"foreignKey:CustomerID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

func (Customer) TableName() string {
	return "customer.customers"
}

type CustomerSession struct {
	ID         uint      `gorm:"primarykey" redis:"id"`
	CustomerID uint      `gorm:"not null;index" redis:"user_id"`
	Token      string    `gorm:"not null;unique;index" redis:"token"`
	ExpiredAt  time.Time `gorm:"not null" redis:"expired_at"`
	IpAddr     string    `gorm:"column:ip_address" redis:"ip_addr"`
	UserAgent  string    `gorm:"not null" redis:"user_agent"`
	CreatedAt  time.Time `gorm:""`
	// Customer   *Customer `gorm:"foreignKey:CustomerID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

func (CustomerSession) TableName() string {
	return "customer.customer_sessions"
}

type CartItem struct {
	CustomerID uint      `gorm:"primaryKey;index"`
	ProductID  uint      `gorm:"primaryKey;index"`
	Quantity   int       `gorm:"not null;default:1"`
	AddedAt    time.Time `gorm:""`
	// Product    *Product  `gorm:"foreignKey:ProductID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

func (CartItem) TableName() string {
	return "customer.customer_carts"
}

type WishlistItem struct {
	CustomerID uint      `gorm:"primaryKey;index"`
	ProductID  uint      `gorm:"primaryKey;index"`
	AddedAt    time.Time `gorm:""`
	// Product    *Product  `gorm:"foreignKey:ProductID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

func (WishlistItem) TableName() string {
	return "customer.customer_wishlists"
}

// Country represents a country in the system
type Country struct {
	ID        uint      `gorm:"primarykey"`
	Name      string    `gorm:"not null;unique;index"`
	Code      string    `gorm:"not null;unique;index;type:varchar(5)"`
	CreatedAt time.Time `gorm:""`
}

func (Country) TableName() string {
	return "backend.countries"
}
