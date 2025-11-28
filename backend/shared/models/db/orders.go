package db

import "time"

type Order struct {
	ID                uint `gorm:"primarykey"`
	CreatedAt         time.Time
	UpdatedAt         time.Time
	UserID            uint    `json:"user_id" gorm:"index;not null"`
	OrderNumber       string  `json:"order_number" gorm:"index;not null;unique"`
	TotalAmount       float64 `json:"total_amount" gorm:"not null"`
	Currency          string  `json:"currency" gorm:"default:'NGN';not null"`
	Status            string  `json:"status" gorm:"default:'pending';not null"`
	BillingAddressID  uint    `json:"billing_address_id" gorm:"index;not null"`
	ShippingAddressID uint    `json:"shipping_address_id" gorm:"index;not null"`
	PaymentID         uint    `json:"payment_id" gorm:"index;not null"`
}

func (Order) TableName() string {
	return "catalog.orders"
}

type OrderReturn struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	OrderID   uint   `json:"order_id" gorm:"index;not null"`
	Reason    string `json:"reason" gorm:"not null"`
}

func (OrderReturn) TableName() string {
	return "catalog.order_returns"
}

type OrderItem struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	OrderID   uint    `json:"order_id" gorm:"index;not null"`
	Name      string  `json:"name" gorm:"not null"`
	Price     float64 `json:"price" gorm:"not null"`
	Quantity  uint    `json:"quantity" gorm:"not null"`
}

func (OrderItem) TableName() string {
	return "catalog.order_items"
}
