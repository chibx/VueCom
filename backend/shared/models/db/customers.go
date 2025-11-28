package db

import "time"

// For timed one time password
type CustomerOTP struct {
	Code       string    `gorm:"not null"`
	ExpiryDate time.Time `gorm:"not null"`
	UserId     uint      `gorm:"column:customer_id;not null"`
}

func (CustomerOTP) TableName() string {
	return "users.customer_otps"
}

type Customer struct {
	ID              uint      `gorm:"primarykey" redis:"id"`
	CreatedAt       time.Time `gorm:"" redis:"created_at"`
	FullName        string    `gorm:"not null;type:varchar(255);index" redis:"full_name"`
	UserName        *string   `gorm:"type:varchar(255);index" redis:"user_name"`
	PasswordHash    *string   `gorm:"" redis:"-"`
	Email           string    `gorm:"unique;not null;type:varchar(255);index" redis:"email"`
	IsEmailVerified bool      `gorm:"default:FALSE;not null" redis:"is_email_verified"`
	PhoneNumber     *string   `gorm:"type:varchar(20)" redis:"phone_number"`
	Image           *string   `gorm:"column:image_url" redis:"image_url"`
	Country         uint      `gorm:"index,not null" redis:"country"`
	City            string    `gorm:"index,not null" redis:"city"`
	State           string    `gorm:"index,not null" redis:"state"`
	ZipCode         string    `gorm:"index,not null" redis:"zip_code"`
	Address         string    `gorm:"index,not null" redis:"address"`
}

func (Customer) TableName() string {
	return "users.customers"
}

type CustomerSession struct {
	ID         uint      `gorm:"primarykey" redis:"id"`
	Token      string    `gorm:"not null" redis:"token"`
	ExpiredAt  time.Time `gorm:"not null" redis:"expired_at"`
	IpAddr     string    `gorm:"column:ip_address" redis:"ip_addr"`
	UserAgent  string    `gorm:"not null" redis:"user_agent"`
	CustomerID uint      `gorm:"not null" redis:"user_id"`
}

func (CustomerSession) TableName() string {
	return "users.customer_sessions"
}
