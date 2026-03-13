package order_service

import (
	"context"
	"sync/atomic"

	"github.com/chibx/vuecom/backend/shared/proto/go/orders"

	"google.golang.org/grpc"
)

type Service struct {
	orders.UnimplementedOrderServiceServer
	nextID atomic.Uint64
	// ...
}

func (s *Service) CreateOrder(ctx context.Context, req *orders.CreateOrderRequest) (*orders.CreateOrderResponse, error) {
	// your logic + possibly publish event
	id := s.nextID.Add(1)
	return &orders.CreateOrderResponse{Id: id}, nil
}

func Register(s *grpc.Server) {
	orders.RegisterOrderServiceServer(s, &Service{})
}
