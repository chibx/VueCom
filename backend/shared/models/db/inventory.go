package db

import "time"

type Inventory struct {
	ID              uint      `gorm:"primarykey" redis:"id"`
	UpdatedAt       time.Time `gorm:"" redis:"updated_at"`
	CreatedAt       time.Time `gorm:"" redis:"created_at"`
	Sku             string    `json:"sku" gorm:"not null;index"`
	ProductId       uint      `json:"product_id" gorm:"not null"`
	WarehouseId     uint      `json:"warehouse_id" gorm:"not null"`
	AvailableQty    int       `json:"available_qty" gorm:"default:0"`
	ReservedQty     int       `json:"reserved_qty" gorm:"default:0"`
	OnHoldQty       int       `json:"on_hold_qty" gorm:"default:0"`
	TotalQty        int       `json:"total_qty" gorm:"default:0"`
	SafetyStock     int       `json:"safety_stock" gorm:"default:0"`
	ReorderLevel    int       `json:"reorder_level" gorm:"default:0"`
	LastRestockedAt time.Time `json:"last_restocked_at" gorm:""`
	LastSoldAt      time.Time `json:"last_sold_at" gorm:""`
}

func (Inventory) TableName() string {
	return "inventory.inventory"
}

type Warehouse struct {
	ID        uint      `gorm:"primarykey" redis:"id"`
	Code      string    `gorm:"not null;index;unique" redis:"code"`
	Name      string    `gorm:"not null;index" redis:"name"`
	Address   string    `gorm:"not null" redis:"address"`
	City      string    `gorm:"not null" redis:"city"`
	StateId   uint      `gorm:"index" redis:"state"`
	CountryId uint      `gorm:"index" redis:"country"`
	IsActive  bool      `gorm:"default:TRUE;not null" redis:"is_active"`
	Capacity  int       `gorm:"default:0" redis:"capacity"`
	CreatedAt time.Time `gorm:"" redis:"created_at"`
	UpdatedAt time.Time `gorm:"" redis:"updated_at"`

	State   *State   `gorm:"foreignKey:StateId;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Country *Country `gorm:"foreignKey:CountryId;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

func (Warehouse) TableName() string {
	return "inventory.warehouses"
}

type StockMovement struct {
	ID           uint      `gorm:"primarykey" redis:"id"`
	UpdatedAt    time.Time `gorm:"" redis:"updated_at"`
	CreatedAt    time.Time `gorm:"" redis:"created_at"`
	InventoryId  uint      `json:"inventory_id" gorm:"not null"`
	Sku          string    `json:"sku" gorm:"not null;index"`
	WarehouseId  uint      `json:"warehouse_id" gorm:"not null"`
	MovementType string    `json:"movement_type" gorm:"not null"`
	Quantity     int       `json:"quantity" gorm:"not null"`
	Reference    string    `json:"reference" gorm:"not null"`
	Notes        string    `json:"notes"`
	CreatedBy    string    `json:"created_by" gorm:"not null"`
}

func (StockMovement) TableName() string {
	return "inventory.stock_movements"
}
