package cache

import (
	"context"
	"time"

	catRes "github.com/chibx/vuecom/backend/services/gateway/api/v1/response/catalog"
	"github.com/chibx/vuecom/backend/services/gateway/internal/cache/keys"
	"github.com/chibx/vuecom/backend/services/gateway/internal/global"
	igrpc "github.com/chibx/vuecom/backend/services/gateway/internal/grpc"
	"github.com/chibx/vuecom/backend/services/gateway/internal/types"
	"github.com/chibx/vuecom/backend/services/gateway/internal/utils"
	"github.com/chibx/vuecom/backend/shared/errors/server"
	"github.com/chibx/vuecom/backend/shared/proto/go/catalog"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

func GetProduct(ctx context.Context, api *types.Api, productId uint32) (*catRes.GetProductResp, error) {
	var product = new(catRes.GetProductResp)
	var err error
	rds := api.Deps.Redis

	err = rds.HGetAll(ctx, keys.ProductKey(productId)).Scan(product)
	// notExist := product.ID == 0
	if err != nil {
		global.Logger().Error("Error getting product data from cache", zap.Error(err))
		resp, err := igrpc.CatalogClient.GetProduct(ctx, &catalog.GetProductRequest{
			Id: uint64(productId),
		})
		product, err = utils.GetProductFromRpc(resp)

		if err != nil {
			if err.Error() == server.ErrDBRecordNotFound.Error() {
				return nil, server.ErrDBRecordNotFound
			}

			global.Logger().Error("Error getting product data from catalog service", zap.Error(err))
			return nil, err
		}

		go SetProduct(ctx, api, product)
	}

	return product, nil
}

func SetProduct(ctx context.Context, api *types.Api, data *catRes.GetProductResp) error {
	rds := api.Deps.Redis
	_, err := rds.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		var err error
		err = pipe.HSet(ctx, keys.ProductKey(data.ID), data /* resp to product */).Err()
		if err != nil {
			return err
		}
		err = pipe.Expire(ctx, keys.ProductKey(data.ID), 10*time.Minute).Err() // Global expiry on the key.
		return err
	})
	if err != nil {
		global.Logger().Error("Error setting product cache", zap.Error(err))
	}
	return err
}
