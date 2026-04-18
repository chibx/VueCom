package inventory

type CreateWarehouseReq struct {
	Code      string `json:"code" validate:"required,min=1,max=10"`
	Name      string `json:"name" validate:"required,min=1,max=100"`
	Address   string `json:"address" validate:"required,min=1,max=255"`
	City      string `json:"city" validate:"required,min=1,max=100"`
	StateID   uint32 `json:"state_id" validate:"omitempty,min=1"`
	CountryID uint32 `json:"country_id" validate:"omitempty,min=1"`
	IsActive  bool   `json:"is_active"`
}

type DeleteWarehouseReq struct {
	WarehouseIDs []uint32 `json:"warehouse_ids"`
}

type CreateStockMovementReq struct {
	InventoryID  uint32 `json:"inventory_id" validate:"required,min=1"`
	SKU          string `json:"sku" validate:"required,min=1,max=50"`
	WarehouseID  uint32 `json:"warehouse_id" validate:"required,min=1"`
	MovementType string `json:"movement_type" validate:"required,oneof=restock sale return adjustment transfer other"`
	Quantity     int32  `json:"quantity" validate:"required"`
	Reference    string `json:"reference" validate:"omitempty,max=100"`
	Notes        string `json:"notes" validate:"omitempty,max=255"`
	CreatedBy    uint32 `json:"created_by" validate:"required,min=1"`
}

type ListStockMovementsReq struct {
	WarehouseID uint32 `json:"warehouse_id" validate:"omitempty,min=1"`
	SKU         string `json:"sku" validate:"omitempty,min=1,max=50"`
}
