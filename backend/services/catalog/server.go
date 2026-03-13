package catalog_service

import (
	"context"
	"sync/atomic"

	"github.com/chibx/vuecom/backend/shared/proto/go/catalog"

	"google.golang.org/grpc"
)

type Service struct {
	catalog.UnimplementedCatalogServiceServer
	nextID atomic.Uint64
	// ...
}

func (s *Service) CreateProduct(ctx context.Context, req *catalog.CreateProductRequest) (*catalog.CreateProductResponse, error) {
	// your logic + possibly publish event
	id := s.nextID.Add(1)
	return &catalog.CreateProductResponse{Id: id}, nil
}

func (s *Service) GetProduct(ctx context.Context, req *catalog.GetProductRequest) (*catalog.GetProductResponse, error) {
	// your logic + possibly publish event
	id := s.nextID.Add(1)
	return &catalog.GetProductResponse{Product: &catalog.Product{Id: id, Name: "Sport Max", Description: "Pro sport", Price: 40000, CategoryId: 1}}, nil
}

func (s *Service) GetCategory(ctx context.Context, req *catalog.GetCategoryRequest) (*catalog.GetCategoryResponse, error) {
	// your logic + possibly publish event
	id := s.nextID.Add(1)
	return &catalog.GetCategoryResponse{Category: &catalog.Category{Id: id, Name: "1234"}}, nil
}

func Register(s *grpc.Server) {
	catalog.RegisterCatalogServiceServer(s, &Service{})
}
