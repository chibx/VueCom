package payment_service

import (
	"context"
	"sync/atomic"

	ordersPr "github.com/chibx/vuecom/backend/shared/proto/go/orders"
)

type Service struct {
	ordersPr.UnimplementedOrderServiceServer
	nextID atomic.Uint64
	// ...
}

func (s *Service) CreateOrder(ctx context.Context, req *ordersPr.CreateOrderRequest) (*ordersPr.CreateOrderResponse, error) {
	// your logic + possibly publish event
	id := s.nextID.Add(1)
	return &ordersPr.CreateOrderResponse{Id: id}, nil
}
