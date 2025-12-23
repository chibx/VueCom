package cache

import (
	"context"
	"errors"
	"time"
	"vuecom/gateway/internal/cache/keys"
	"vuecom/gateway/internal/types"
	dbModels "vuecom/shared/models/db"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func GetAppData(ctx context.Context, api *types.Api) (*dbModels.AppData, error) {
	db := api.Deps.DB
	cache := api.Deps.Redis
	appData := new(dbModels.AppData)

	err := cache.Get(ctx, keys.APP_DATA_KEY).Scan(&appData)
	if err != nil {
		if !errors.Is(err, redis.Nil) {
			return nil, fiber.NewError(fiber.StatusInternalServerError, "Error fetching app data")
		}

		// Key Not Found
		if err := db.WithContext(ctx).Limit(1).First(appData).Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fiber.NewError(fiber.StatusInternalServerError, "Error fetching app data")
		}

		if appData == nil {
			return nil, nil
		}

		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
			defer cancel()
			// Cache for 5 minutes
			// TODO: The error could be sent to a logger or monitoring tool/service
			_ = cache.Set(ctx, keys.APP_DATA_KEY, appData, time.Minute*5).Err()
		}()

	}

	return appData, nil
}
