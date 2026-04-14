package inventory_service

import (
	"context"
	"sync/atomic"

	inventoryPr "github.com/chibx/vuecom/backend/shared/proto/go/inventory"
)

type Service struct {
	inventoryPr.UnimplementedInventoryServiceServer
	nextID atomic.Uint64
	// ...
}

func (s *Service) CreateProductRecord(ctx context.Context, req *inventoryPr.AddProductRequest) (*inventoryPr.AddProductResponse, error) {
	// your logic + possibly publish event
	// id := s.nextID.Add(1)
	return &inventoryPr.AddProductResponse{}, nil
}
