package inventory

import "time"

// CREATE TYPE stock_movement_type AS ENUM('restock', 'sale', 'return', 'adjustment', 'transfer', 'other')
type StockMovementType string

const (
	MovementRestock    StockMovementType = "restock"
	MovementSale       StockMovementType = "sale"
	MovementReturn     StockMovementType = "return"
	MovementAdjustment StockMovementType = "adjustment"
	MovementTransfer   StockMovementType = "transfer"
	MovementOther      StockMovementType = "other"
)

type Inventory struct {
	ID              uint      `gorm:"primarykey" redis:"id"`
	UpdatedAt       time.Time `gorm:"" redis:"updated_at"`
	CreatedAt       time.Time `gorm:"" redis:"created_at"`
	SKU             string    `json:"sku" gorm:"not null;index" redis:"sku"`
	ProductId       uint      `json:"product_id" gorm:"not null" redis:"product_id"`
	WarehouseId     uint      `json:"warehouse_id" gorm:"not null" redis:"warehouse_id"`
	AvailableQty    int       `json:"available_qty" gorm:"default:0" redis:"available_qty"`
	ReservedQty     int       `json:"reserved_qty" gorm:"default:0" redis:"reserved_qty"`
	OnHoldQty       int       `json:"on_hold_qty" gorm:"default:0" redis:"on_hold_qty"`
	TotalQty        int       `json:"total_qty" gorm:"<-:false" redis:"total_qty"`
	SafetyStock     int       `json:"safety_stock" gorm:"default:0" redis:"safety_stock"`
	ReorderLevel    int       `json:"reorder_level" gorm:"default:0" redis:"reorder_level"`
	LastRestockedAt time.Time `json:"last_restocked_at" gorm:"" redis:"last_restocked_at"`
	LastSoldAt      time.Time `json:"last_sold_at" gorm:"" redis:"last_sold_at"`
}

type Warehouse struct {
	ID        uint      `gorm:"primarykey" redis:"id"`
	Code      string    `gorm:"not null;index;unique" redis:"code"`
	Name      string    `gorm:"not null;index" redis:"name"`
	Address   *string   `gorm:"not null" redis:"address"`
	City      *string   `gorm:"not null" redis:"city"`
	StateID   *uint     `gorm:"index" redis:"state"`
	CountryID *uint     `gorm:"index" redis:"country"`
	IsActive  bool      `gorm:"default:TRUE;not null" redis:"is_active"`
	Capacity  int       `gorm:"default:0" redis:"capacity"`
	CreatedAt time.Time `gorm:"" redis:"created_at"`
	UpdatedAt time.Time `gorm:"" redis:"updated_at"`

	// State   *State   `gorm:"foreignKey:StateId;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" redis:"-"`
	// Country *Country `gorm:"foreignKey:CountryId;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" redis:"-"`
}

type StockMovement struct {
	ID           uint              `gorm:"primarykey" redis:"id"`
	SKU          string            `json:"sku" gorm:"not null;index" redis:"sku"`
	InventoryId  uint              `json:"inventory_id" gorm:"not null" redis:"inventory_id"`
	WarehouseId  uint              `json:"warehouse_id" gorm:"not null" redis:"warehouse_id"`
	MovementType StockMovementType `json:"movement_type" gorm:"not null" redis:"movement_type"`
	Quantity     int               `json:"quantity" gorm:"not null" redis:"quantity"`
	Reference    string            `json:"reference" gorm:"not null" redis:"reference"`
	Notes        string            `json:"notes" redis:"notes"`
	CreatedBy    uint              `json:"created_by" gorm:"not null" redis:"created_by"`
	CreatedAt    time.Time         `gorm:"" redis:"created_at"`
	// UpdatedAt    time.Time `gorm:"" redis:"updated_at"`
}
