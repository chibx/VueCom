package models

import "time"

type CreateAppData struct {
	Name       string `json:"app_name"`
	AdminRoute string `json:"-"`
	Plan       int    `json:"app_plan"`
	LogoUrl    string `json:"app_logo"`
}

type OnlyID struct {
	ID int `json:"id"`
}

// JWT Format sent back to the client dashboard
type BackendJWTPayload struct {
	UserId int    `json:"user_id"`
	Role   string `json:"role"`
	Exp    int    `json:"exp"`
}

// JWT Format sent back to the client
type CustomerJWTPayload struct {
	UserId int    `json:"user_id"`
	Role   string `json:"role"`
	Exp    int    `json:"exp"`
}

type sharedUserProps struct {
	FullName        string  `json:"full_name" gorm:"not null;type:varchar(255);index"`
	UserName        *string `json:"user_name" gorm:"type:varchar(255);index"`
	Email           string  `json:"email" gorm:"unique;not null;type:varchar(255);index"`
	PhoneNumber     *string `json:"phone_number" gorm:"type:varchar(20)"`
	Image           *string `json:"image"`
	Country         uint    `json:"country" gorm:"index"`
	IsEmailVerified bool    `json:"email_verified" gorm:"default:FALSE;not null"`
}

type ApiCustomer struct {
	sharedUserProps
}

// Base Backend Panel User
type ApiBackendUser struct {
	sharedUserProps
	Role string `json:"role" gorm:"type:varchar(50)"`
}

type ApiProducts struct {
	Name string `json:"name" gorm:"not null;index;type:text"`
	SKU  string `json:"sku" gorm:"not null;index"`
	// Just made the precision to be 15 (Don't know how bad the ecomomy of some countries are)
	Price      float64   `json:"price" gorm:"not null;type:numeric(15, 2)"`
	DscPercent float64   `json:"dsc_percent" gorm:"type:numeric(5, 2)"`
	DscPeriod  time.Time `json:"dsc_period" gorm:""`
	Enabled    bool      `json:"enabled" gorm:""`
	// Warranty   string    `json:"warranty"`
	Description string
	Url         string `json:"url"`
}

type CustomerReviews struct {
	Id        uint      `gorm:"primary"`
	Text      string    `gorm:"text"`
	Rating    int8      `gorm:"rating"`
	CreatedAt time.Time ``
	EditTimes int8      `gorm:""`
	UserId    uint      `gorm:""`
	ProductId uint      `gorm:""`
}

type f struct {
}
