package inventory

import "time"

type Inventory struct {
	ID              uint      `gorm:"primarykey" redis:"id"`
	UpdatedAt       time.Time `gorm:"" redis:"updated_at"`
	CreatedAt       time.Time `gorm:"" redis:"created_at"`
	Sku             string    `json:"sku" gorm:"not null;index" redis:"sku"`
	ProductId       uint      `json:"product_id" gorm:"not null" redis:"product_id"`
	WarehouseId     uint      `json:"warehouse_id" gorm:"not null" redis:"warehouse_id"`
	AvailableQty    int       `json:"available_qty" gorm:"default:0" redis:"available_qty"`
	ReservedQty     int       `json:"reserved_qty" gorm:"default:0" redis:"reserved_qty"`
	OnHoldQty       int       `json:"on_hold_qty" gorm:"default:0" redis:"on_hold_qty"`
	TotalQty        int       `json:"total_qty" gorm:"default:0" redis:"total_qty"`
	SafetyStock     int       `json:"safety_stock" gorm:"default:0" redis:"safety_stock"`
	ReorderLevel    int       `json:"reorder_level" gorm:"default:0" redis:"reorder_level"`
	LastRestockedAt time.Time `json:"last_restocked_at" gorm:"" redis:"last_restocked_at"`
	LastSoldAt      time.Time `json:"last_sold_at" gorm:"" redis:"last_sold_at"`
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

	// State   *State   `gorm:"foreignKey:StateId;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" redis:"-"`
	// Country *Country `gorm:"foreignKey:CountryId;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" redis:"-"`
}

func (Warehouse) TableName() string {
	return "inventory.warehouses"
}

type StockMovement struct {
	ID           uint      `gorm:"primarykey" redis:"id"`
	UpdatedAt    time.Time `gorm:"" redis:"updated_at"`
	CreatedAt    time.Time `gorm:"" redis:"created_at"`
	InventoryId  uint      `json:"inventory_id" gorm:"not null" redis:"inventory_id"`
	Sku          string    `json:"sku" gorm:"not null;index" redis:"sku"`
	WarehouseId  uint      `json:"warehouse_id" gorm:"not null" redis:"warehouse_id"`
	MovementType string    `json:"movement_type" gorm:"not null" redis:"movement_type"`
	Quantity     int       `json:"quantity" gorm:"not null" redis:"quantity"`
	Reference    string    `json:"reference" gorm:"not null" redis:"reference"`
	Notes        string    `json:"notes" redis:"notes"`
	CreatedBy    string    `json:"created_by" gorm:"not null" redis:"created_by"`
}

func (StockMovement) TableName() string {
	return "inventory.stock_movements"
}
