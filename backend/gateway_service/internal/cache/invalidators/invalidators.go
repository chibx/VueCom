package invalidators

import (
	"context"

	"github.com/redis/go-redis/v9"
)

const DEFAULT_BATCH_SIZE = 100
const USER_KEY_PATTERN = "user:*"
const PRODUCT_KEY_PATTERN = "product:*"
const BACKEND_USER_KEY_PATTERN = "backend_user:*"

// TODO: I believe this simple function could be better optimized later

func InvalidateCache(rdb *redis.Client, pattern string) {
	go func() {
		var cursor uint64
		var deletedCount int
		ctx := context.Background()
		for {
			keys, nextCursor, _ := rdb.Scan(ctx, cursor, pattern, DEFAULT_BATCH_SIZE).Result()
			// if err != nil {
			// 	return deletedCount, err
			// }
			if len(keys) > 0 {
				_, _ = rdb.Unlink(ctx, keys...).Result()
				// if err != nil {
				// 	return deletedCount, err
				// }
				deletedCount += len(keys)
			}
			cursor = nextCursor
			if cursor == 0 {
				break
			}
		}
	}()
	// return deletedCount, nil
}
