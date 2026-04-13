package cache

import (
	"context"
	"time"

	"github.com/chibx/vuecom/backend/services/gateway/internal/cache/keys"
	"github.com/chibx/vuecom/backend/services/gateway/internal/global"
	igrpc "github.com/chibx/vuecom/backend/services/gateway/internal/grpc"
	"github.com/chibx/vuecom/backend/services/gateway/internal/types"
	"github.com/chibx/vuecom/backend/shared/errors/server"
	catModels "github.com/chibx/vuecom/backend/shared/models/db/catalog"
	"github.com/chibx/vuecom/backend/shared/proto/go/catalog"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

func GetProduct(ctx context.Context, api types.Api, productId uint32) (*catModels.Product, error) {
	var product = new(catModels.Product)
	var err error
	rds := api.Deps.Redis

	err = rds.HGetAll(ctx, keys.ProductKey(productId)).Scan(product)
	// notExist := product.ID == 0
	if err != nil {
		global.Logger().Error("Error getting product data from cache", zap.Error(err))
		resp, err := igrpc.CatalogClient.GetProduct(ctx, &catalog.GetProductRequest{
			Id: uint64(productId),
		})
		if err != nil {
			if err.Error() == server.ErrDBRecordNotFound.Error() {
				return nil, server.ErrDBRecordNotFound
			}

			global.Logger().Error("Error getting product data from catalog service", zap.Error(err))
			return nil, err
		}

		go func() {
			_, err := rds.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
				var err error
				err = pipe.HSet(ctx, keys.ProductKey(productId), resp /* resp to product */).Err()
				if err != nil {
					return err
				}
				err = pipe.Expire(ctx, keys.ProductKey(productId), 10*time.Minute).Err() // Global expiry on the key.
				return err
			})
			if err != nil {
				global.Logger().Error("Error setting product cache", zap.Error(err))
			}
		}()
	}

	return product, nil
}
