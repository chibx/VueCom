package inventory

type CreateWarehouseResp struct {
	ID uint32 `json:"id"`
}

type CreateStockMovementResp struct {
	ID uint32 `json:"id"`
}

type ListStockMovementsResp struct {
	StockMovements []StockMovement `json:"stock_movements"`
}

type StockMovement struct {
	ID           uint32 `json:"id"`
	InventoryID  uint32 `json:"inventory_id"`
	SKU          string `json:"sku"`
	WarehouseID  uint32 `json:"warehouse_id"`
	MovementType string `json:"movement_type"`
	Quantity     int32  `json:"quantity"`
	Reference    string `json:"reference"`
	Notes        string `json:"notes"`
	CreatedBy    uint32 `json:"created_by"`
	CreatedAt    string `json:"created_at"`
}
