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

type CustomerAddress struct {
	ID            uint     `gorm:"primarykey" redis:"id"`
	CustomerID    uint     `gorm:"not null;index" redis:"user_id"`
	StreetAddress string   `gorm:"index;not null" redis:"address"`
	City          string   `gorm:"index;not null" redis:"city"`
	ZipCode       string   `gorm:"index;not null" redis:"zip_code"`
	StateID       *uint    `gorm:"column:state;index" redis:"state"`
	CountryID     *uint    `gorm:"column:country;index" redis:"country"`
	State         *State   `gorm:"foreignKey:StateID;"`
	Country       *Country `gorm:"foreignKey:CountryID;"`
}

type Customer struct {
	ID              uint              `gorm:"primarykey" redis:"id"`
	CreatedAt       time.Time         `gorm:"" redis:"created_at"`
	FullName        string            `gorm:"not null;type:varchar(255);index" redis:"full_name"`
	UserName        *string           `gorm:"type:varchar(255);index" redis:"user_name"`
	PasswordHash    *string           `gorm:"" redis:"-"`
	Email           string            `gorm:"unique;not null;type:varchar(255);index" redis:"email"`
	IsEmailVerified bool              `gorm:"default:FALSE;not null" redis:"is_email_verified"`
	PhoneNumber     *string           `gorm:"type:varchar(20);" redis:"phone_number"`
	Image           *string           `gorm:"column:image_url" redis:"image_url"`
	BillingAddress  *CustomerAddress  `gorm:"foreignKey:CustomerID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" redis:"-"`
	ShippingAddress *CustomerAddress  `gorm:"foreignKey:CustomerID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" redis:"-"`
	Sessions        []CustomerSession `gorm:"foreignKey:CustomerID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" redis:"-"`
	OTPs            []CustomerOTP     `gorm:"foreignKey:CustomerID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" redis:"-"`
	Wishlist        []WishlistItem    `gorm:"foreignKey:CustomerID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" redis:"-"`
	Cart            []CartItem        `gorm:"foreignKey:CustomerID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" redis:"-"`
	// Wishlist        []Product         `gorm:"many2many:customer.customer_wishlists;"`
}

func (Customer) TableName() string {
	return "customer.customers"
}

type CustomerSession struct {
	CustomerID uint      `gorm:"not null" redis:"user_id"`
	Token      string    `gorm:"not null" redis:"-"` // redis key would be the token
	ExpiredAt  time.Time `gorm:"not null" redis:"expired_at"`
	IpAddr     string    `gorm:"column:ip_address" redis:"ip_address"`
	UserAgent  string    `gorm:"not null" redis:"user_agent"`
	Customer   *Customer `gorm:"foreignKey:CustomerID;constraint:OnDelete:CASCADE;" redis:"-"`
}

func (CustomerSession) TableName() string {
	return "customer.customer_sessions"
}

type CartItem struct {
	CustomerID uint      `gorm:"primaryKey;index"`
	ProductID  uint      `gorm:"primaryKey;index"`
	Quantity   int       `gorm:"not null;default:1"`
	AddedAt    time.Time `gorm:""`
	Product    *Product  `gorm:"foreignKey:ProductID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

func (CartItem) TableName() string {
	return "customer.customer_carts"
}

type WishlistItem struct {
	CustomerID uint      `gorm:"primaryKey;index"`
	ProductID  uint      `gorm:"primaryKey;index"`
	AddedAt    time.Time `gorm:""`
	Product    *Product  `gorm:"foreignKey:ProductID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

func (WishlistItem) TableName() string {
	return "customer.customer_wishlists"
}
