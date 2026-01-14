package orders

import "time"

type Order struct {
	ID                uint      `gorm:"primarykey" redis:"id"`
	CreatedAt         time.Time `gorm:"" redis:"created_at"`
	UpdatedAt         time.Time `gorm:"" redis:"updated_at"`
	UserID            uint      `json:"user_id" gorm:"index;not null" redis:"user_id"`
	OrderNumber       string    `json:"order_number" gorm:"index;not null;unique" redis:"order_number"`
	TotalAmount       float64   `json:"total_amount" gorm:"not null" redis:"total_amount"`
	Currency          string    `json:"currency" gorm:"default:'NGN';not null" redis:"currency"`
	Status            string    `json:"status" gorm:"default:'pending';not null" redis:"status"`
	BillingAddressID  uint      `json:"billing_address_id" gorm:"index;not null" redis:"billing_address_id"`
	ShippingAddressID uint      `json:"shipping_address_id" gorm:"index;not null" redis:"shipping_address_id"`
	PaymentID         uint      `json:"payment_id" gorm:"index;not null" redis:"payment_id"`
}

func (Order) TableName() string {
	return "catalog.orders"
}

type OrderReturn struct {
	ID        uint      `gorm:"primarykey" redis:"id"`
	CreatedAt time.Time `gorm:"" redis:"created_at"`
	OrderID   uint      `json:"order_id" gorm:"index;not null" redis:"order_id"`
	Reason    string    `json:"reason" gorm:"not null" redis:"reason"`
}

func (OrderReturn) TableName() string {
	return "catalog.order_returns"
}

type OrderItem struct {
	ID        uint      `gorm:"primarykey" redis:"id"`
	CreatedAt time.Time `gorm:"" redis:"created_at"`
	UpdatedAt time.Time `gorm:"" redis:"updated_at"`
	OrderID   uint      `json:"order_id" gorm:"index;not null" redis:"order_id"`
	Name      string    `json:"name" gorm:"not null" redis:"name"`
	Price     float64   `json:"price" gorm:"not null" redis:"price"`
	Quantity  uint      `json:"quantity" gorm:"not null" redis:"quantity"`
}

func (OrderItem) TableName() string {
	return "catalog.order_items"
}
