package models

import "time"

type CreateAppData struct {
	Name       string `json:"appName"`
	AdminRoute string `json:"adminRoute"`
	Plan       int    `json:"appPlan"`
	ImageUrl   string `json:"appImageUrl"`
}

type OnlyID struct {
	ID int `json:"id"`
}

// JWT Format sent back to the client dashboard
type BackendJWTPayload struct {
	UserId int    `json:"userId"`
	Role   string `json:"role"`
	Exp    int    `json:"exp"`
}

// JWT Format sent back to the client
type CustomerJWTPayload struct {
	UserId int    `json:"userId"`
	Role   string `json:"role"`
	Exp    int    `json:"exp"`
}

// Base Backend Panel User
type ApiUser struct {
	FullName        string  `json:"fullName" gorm:"not null;type:varchar(255);index"`
	UserName        *string `json:"userName" gorm:"type:varchar(255);index"`
	Email           string  `json:"email" gorm:"unique;not null;type:varchar(255);index"`
	IsEmailVerified bool    `json:"emailVerified" gorm:"default:FALSE;not null"`
	PhoneNumber     *string `json:"phoneNumber" gorm:"type:varchar(20)"`
	Image           *string `json:"image"`
	Country         int     `json:"country" gorm:"index"`
	Role            string  `json:"role" gorm:"type:varchar(50)"`
}

type ApiProducts struct {
	Name       string    `json:"name" gorm:"not null;index;type:text"`
	SKU        string    `json:"sku" gorm:"not null;index"`
	Price      string    `json:"price" gorm:"not null;"`
	DscPercent string    `json:"dscPercent"`
	DscPeriod  time.Time `json:"dscPeriod" gorm:""`
	Enabled    bool      `json:"enabled" gorm:""`
	Warranty   string    `json:"warranty"`
}
