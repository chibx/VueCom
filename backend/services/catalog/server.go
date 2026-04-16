package catalog_service

import (
	"context"
	"sync/atomic"

	"github.com/chibx/vuecom/backend/services/catalog/internal/db"
	"github.com/chibx/vuecom/backend/services/catalog/internal/global"
	igrpc "github.com/chibx/vuecom/backend/services/catalog/internal/grpc"
	"github.com/chibx/vuecom/backend/services/catalog/internal/utils"
	"github.com/chibx/vuecom/backend/shared/events"
	catalogPr "github.com/chibx/vuecom/backend/shared/proto/go/catalog"
	"github.com/chibx/vuecom/backend/shared/proto/go/inventory"
	"go.uber.org/zap"
)

type Service struct {
	catalogPr.UnimplementedCatalogServiceServer
	nextID atomic.Uint32
	// ...
}

func (s *Service) CreateProduct(ctx context.Context, req *catalogPr.CreateProductRequest) (*catalogPr.CreateProductResponse, error) {
	product := utils.CreateProdRpcToDBModel(req)

	err := global.Repo.RunInTx(func(c *db.CatalogDB) error {
		err := c.CreateProduct(ctx, product)
		if err != nil {
			global.Logger.Error("Failed to create product [Tx 1]", zap.Error(err))
			return err
		}
		err = c.CreateProductToCategory(ctx, product.ID, req.PresetValues)
		if err != nil {
			global.Logger.Error("Failed to insert product category relationship [Tx 2]", zap.Error(err))
			return err
		}

		err = c.CreateProductRelation(
			ctx,
			product.ID,
			req.RelatedProducts,
			req.UpSellProducts,
			req.CrossSell,
		)
		if err != nil {
			global.Logger.Error("Failed to insert product to product relationship [Tx 3]", zap.Error(err))
			return err
		}

		return nil
	})
	if err != nil {
		global.Logger.Error("Failed to create product", zap.Error(err))
		return nil, err
	}

	// err = pubsub.DefPubSub.Publish(events.INVENTORY_QUEUE, string(events.PRODUCT_CREATION), pubTypes.CreateInventoryReq{
	// 	ProductId: product.ID,
	// 	Quantity:  req.Quantity,
	// })

	var inventoryWarehouses = make([]*inventory.WarehouseInfo, 0, len(req.Warehouses))
	for _, v := range req.Warehouses {
		inventoryWarehouses = append(inventoryWarehouses, &inventory.WarehouseInfo{
			Id: v.Id, Quantity: v.Quantity,
		})
	}

	_, err = igrpc.InventoryClient.CreateProductRecord(ctx, &inventory.AddProductRequest{
		ProductId:     product.ID,
		WarehouseInfo: inventoryWarehouses,
	})

	if err != nil {
		// global.Logger.Error("Failed to publish event", zap.Error(err), zap.String("event", string(events.PRODUCT_CREATION)))
		global.Logger.Error("Failed to create inventory record", zap.Error(err), zap.String("event", string(events.PRODUCT_CREATION)))
		return nil, err
	}
	return &catalogPr.CreateProductResponse{Id: product.ID}, nil
}

func (s *Service) GetProduct(ctx context.Context, req *catalogPr.GetProductRequest) (*catalogPr.GetProductResponse, error) {
	// your logic + possibly publish event
	id := s.nextID.Add(1)
	return &catalogPr.GetProductResponse{Id: id}, nil
}

func (s *Service) GetCategory(ctx context.Context, req *catalogPr.GetCategoryRequest) (*catalogPr.GetCategoryResponse, error) {
	// your logic + possibly publish event
	id := s.nextID.Add(1)
	return &catalogPr.GetCategoryResponse{Category: &catalogPr.Category{Id: uint64(id), Name: "1234"}}, nil
}
