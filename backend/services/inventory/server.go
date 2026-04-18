package inventory_service

import (
	"context"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/chibx/vuecom/backend/services/inventory/internal/global"
	"github.com/chibx/vuecom/backend/shared/models/db/inventory"
	inventoryPr "github.com/chibx/vuecom/backend/shared/proto/go/inventory"
	"go.uber.org/zap"
)

type Service struct {
	inventoryPr.UnimplementedInventoryServiceServer
	nextID atomic.Uint64
	// ...
}

func (s *Service) CreateProductRecord(ctx context.Context, req *inventoryPr.AddProductRequest) (*inventoryPr.AddProductResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("request cannot be nil")
	}

	if req.ProductId == 0 {
		return nil, fmt.Errorf("product_id must be set")
	}

	if len(req.WarehouseInfo) == 0 {
		return nil, fmt.Errorf("warehouse_info must contain at least one warehouse")
	}

	items := make([]*inventory.Inventory, 0, len(req.WarehouseInfo))
	for _, wh := range req.WarehouseInfo {
		if wh == nil {
			continue
		}
		if wh.Id == 0 {
			return nil, fmt.Errorf("warehouse id must be set for each warehouse_info item")
		}

		item := &inventory.Inventory{
			SKU:          req.Sku,
			ProductId:    uint(req.ProductId),
			WarehouseId:  uint(wh.Id),
			AvailableQty: int(wh.Quantity),
		}
		if wh.Quantity > 0 {
			t := time.Now()
			item.LastRestockedAt = &t
		}
		items = append(items, item)
	}

	if len(items) == 0 {
		return nil, fmt.Errorf("no valid warehouse inventory items to create")
	}

	err := global.Repo.CreateInventoryRecords(ctx, items)
	if err != nil {
		global.Logger.Error("failed to create inventory records", zap.Error(err), zap.Uint32("product_id", req.ProductId))
		return nil, err
	}

	return &inventoryPr.AddProductResponse{}, nil
}

func (s *Service) HasAnyWarehouse(ctx context.Context, req *inventoryPr.WarehouseExistReq) (*inventoryPr.WarehouseExistResp, error) {
	exists, err := global.Repo.HasAnyWarehouse(ctx)
	if err != nil {
		global.Logger.Error("failed to check warehouse existence", zap.Error(err))
		return nil, err
	}

	return &inventoryPr.WarehouseExistResp{Exists: exists}, nil
}

func (s *Service) ListWarehouses(ctx context.Context, req *inventoryPr.ListWarehousesReq) (*inventoryPr.ListWarehousesResp, error) {
	warehouses, err := global.Repo.ListWarehouses(ctx)
	if err != nil {
		global.Logger.Error("failed to list warehouses", zap.Error(err))
		return nil, err
	}

	resp := &inventoryPr.ListWarehousesResp{}
	for _, wh := range warehouses {

		resp.Warehouses = append(resp.Warehouses, &inventoryPr.Warehouse{
			Id:        uint32(wh.ID),
			Code:      wh.Code,
			Name:      wh.Name,
			Address:   wh.Address,
			City:      wh.City,
			StateId:   uint32(wh.StateID),
			CountryId: uint32(wh.CountryID),
			IsActive:  wh.IsActive,
			CreatedAt: wh.CreatedAt.Format(time.RFC3339),
			UpdatedAt: wh.UpdatedAt.Format(time.RFC3339),
		})
	}

	return resp, nil
}

func (s *Service) CreateWarehouse(ctx context.Context, req *inventoryPr.CreateWarehouseReq) (*inventoryPr.CreateWarehouseResp, error) {
	warehouse := &inventory.Warehouse{
		Code:      req.Code,
		Name:      req.Name,
		Address:   req.Address,
		City:      req.City,
		StateID:   uint(req.StateId),
		CountryID: uint(req.CountryId),
		IsActive:  req.IsActive,
	}

	err := global.Repo.CreateWarehouse(ctx, warehouse)
	if err != nil {
		global.Logger.Error("failed to create warehouse", zap.Error(err))
		return nil, err
	}

	return &inventoryPr.CreateWarehouseResp{Id: uint32(warehouse.ID)}, nil
}

func (s *Service) DeleteWarehouse(ctx context.Context, req *inventoryPr.DeleteWarehouseReq) (*inventoryPr.DeleteWarehouseResp, error) {
	err := global.Repo.DeleteWarehouse(ctx, req.WarehouseIds)
	if err != nil {
		global.Logger.Error("failed to delete warehouse(s)", zap.Error(err), zap.Uint32s("ids", req.WarehouseIds))
		return nil, err
	}

	return &inventoryPr.DeleteWarehouseResp{}, nil
}

func (s *Service) CreateStockMovement(ctx context.Context, req *inventoryPr.CreateStockMovementReq) (*inventoryPr.CreateStockMovementResp, error) {
	movement := &inventory.StockMovement{
		InventoryId:  uint(req.InventoryId),
		SKU:          req.Sku,
		WarehouseId:  uint(req.WarehouseId),
		MovementType: inventory.StockMovementType(req.MovementType),
		Quantity:     int(req.Quantity),
		Reference:    req.Reference,
		Notes:        req.Notes,
		CreatedBy:    uint(req.CreatedBy),
	}

	err := global.Repo.CreateStockMovement(ctx, movement)
	if err != nil {
		global.Logger.Error("failed to create stock movement", zap.Error(err))
		return nil, err
	}

	return &inventoryPr.CreateStockMovementResp{Id: uint32(movement.ID)}, nil
}

func (s *Service) ListStockMovements(ctx context.Context, req *inventoryPr.ListStockMovementsReq) (*inventoryPr.ListStockMovementsResp, error) {
	movements, err := global.Repo.ListStockMovements(ctx, uint(req.WarehouseId), req.Sku)
	if err != nil {
		global.Logger.Error("failed to list stock movements", zap.Error(err))
		return nil, err
	}

	resp := &inventoryPr.ListStockMovementsResp{}
	for _, m := range movements {
		resp.StockMovements = append(resp.StockMovements, &inventoryPr.StockMovement{
			Id:           uint32(m.ID),
			InventoryId:  uint32(m.InventoryId),
			Sku:          m.SKU,
			WarehouseId:  uint32(m.WarehouseId),
			MovementType: string(m.MovementType),
			Quantity:     int32(m.Quantity),
			Reference:    m.Reference,
			Notes:        m.Notes,
			CreatedBy:    uint32(m.CreatedBy),
			CreatedAt:    m.CreatedAt.Format(time.RFC3339),
		})
	}

	return resp, nil
}
