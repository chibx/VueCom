package invalidators

import (
	"context"

	"github.com/chibx/vuecom/backend/services/gateway/internal/types"

	"go.uber.org/zap"
)

const DEFAULT_BATCH_SIZE = 100
const USER_KEY_PATTERN = "user:*"
const PRODUCT_KEY_PATTERN = "product:*"
const BACKEND_USER_KEY_PATTERN = "backend_user:*"

// TODO: I believe this simple function could be better optimized later

func InvalidateCache(ctx context.Context, api *types.Api, pattern string) {
	logger := api.Deps.Logger
	rdb := api.Deps.Redis
	go func() {
		var cursor uint64
		var deletedCount int
		for {
			keys, nextCursor, err := rdb.Scan(ctx, cursor, pattern, DEFAULT_BATCH_SIZE).Result()
			logger.Error("Error getting pattern "+pattern, zap.Error(err))
			// if err != nil {
			// 	return deletedCount, err
			// }
			if len(keys) > 0 {
				_, err = rdb.Unlink(ctx, keys...).Result()
				if err != nil {
					logger.Error("Error unlinking cache keys", zap.Error(err))
				}
				deletedCount += len(keys)
			}
			cursor = nextCursor
			if cursor == 0 {
				break
			}
		}

		logger.Info("Total Keys Invalidated for pattern "+pattern, zap.Int("keys", deletedCount))
	}()
	// return deletedCount, nil
}
