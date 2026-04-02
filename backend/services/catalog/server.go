package catalog_service

import (
	"context"
	"sync/atomic"

	"github.com/chibx/vuecom/backend/services/catalog/internal/global"
	"github.com/chibx/vuecom/backend/services/catalog/internal/utils"
	catalogPr "github.com/chibx/vuecom/backend/shared/proto/go/catalog"
	"go.uber.org/zap"
)

type Service struct {
	catalogPr.UnimplementedCatalogServiceServer
	nextID atomic.Uint64
	// ...
}

func (s *Service) CreateProduct(ctx context.Context, req *catalogPr.CreateProductRequest) (*catalogPr.CreateProductResponse, error) {
	product := utils.CreateProdRpcToDBModel(req)

	err := global.Repo.CreateProduct(ctx, product)
	if err != nil {
		global.Logger.Error("Failed to create product", zap.Error(err))
		return nil, err
	}
	return &catalogPr.CreateProductResponse{Id: product.ID}, nil
}

func (s *Service) GetProduct(ctx context.Context, req *catalogPr.GetProductRequest) (*catalogPr.GetProductResponse, error) {
	// your logic + possibly publish event
	id := s.nextID.Add(1)
	return &catalogPr.GetProductResponse{Product: &catalogPr.Product{Id: id, Name: "Sport Max", Description: "Pro sport", Price: 40000, CategoryId: 1}}, nil
}

func (s *Service) GetCategory(ctx context.Context, req *catalogPr.GetCategoryRequest) (*catalogPr.GetCategoryResponse, error) {
	// your logic + possibly publish event
	id := s.nextID.Add(1)
	return &catalogPr.GetCategoryResponse{Category: &catalogPr.Category{Id: id, Name: "1234"}}, nil
}
