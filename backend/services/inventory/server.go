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
