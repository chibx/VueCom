package cache

import (
	"context"
	"errors"
	"time"
	"vuecom/gateway/internal/cache/keys"
	"vuecom/gateway/internal/types"
	"vuecom/shared/errors/server"
	dbModels "vuecom/shared/models/db"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func GetAppData(ctx context.Context, api *types.Api) (*dbModels.AppData, error) {
	db := api.Deps.DB
	cache := api.Deps.Redis
	logger := api.Deps.Logger
	appData := new(dbModels.AppData)

	err := cache.Get(ctx, keys.APP_DATA_KEY).Scan(&appData)
	if err != nil {
		if !errors.Is(err, redis.Nil) {
			logger.Error("Error fetching app data from redis", zap.Error(err))
			return nil, server.NewServerErr(fiber.StatusInternalServerError, "Error fetching app data")
		}

		// Key Not Found
		// err := db.WithContext(ctx).Limit(1).First(appData).Error;
		appData, err = db.AppData().GetAppData(ctx)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, server.NewServerErr(fiber.StatusNotFound, "App data not found")
			}

			logger.Error("Error fetching app data from database", zap.Error(err))
			return nil, server.NewServerErr(fiber.StatusInternalServerError, "Error fetching app data")
		}

		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
			defer cancel()
			// Cache for 5 minutes
			// TODO: The error could be sent to a logger or monitoring tool/service
			err = cache.Set(ctx, keys.APP_DATA_KEY, appData, time.Minute*5).Err()
			if err != nil {
				logger.Error("Error setting appdata cache")
			}
		}()

	}

	return appData, nil
}
