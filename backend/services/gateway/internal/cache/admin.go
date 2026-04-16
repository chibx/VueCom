package cache

import (
	"context"
	"errors"
	"time"

	appModels "github.com/chibx/vuecom/backend/shared/models/db/appdata"

	"github.com/chibx/vuecom/backend/services/gateway/internal/cache/keys"
	"github.com/chibx/vuecom/backend/services/gateway/internal/global"
	"github.com/chibx/vuecom/backend/services/gateway/internal/types"

	// userModels "github.com/chibx/vuecom/backend/shared/models/db/users"

	serverErrors "github.com/chibx/vuecom/backend/shared/errors/server"
	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

func GetAppData(ctx context.Context, api *types.Api) (*appModels.AppData, error) {
	db := api.Deps.DB
	cache := api.Deps.Redis
	logger := global.Logger
	appData := new(appModels.AppData)

	err := cache.HGetAll(ctx, keys.APP_DATA_KEY).Scan(&appData)
	if err != nil {
		if !errors.Is(err, redis.Nil) {
			logger.Error("Error fetching app data from redis", zap.Error(err))
			return nil, serverErrors.NewServerErr(fiber.StatusInternalServerError, "Error fetching app data")
		}

		// Key Not Found
		// err := db.WithContext(ctx).Limit(1).First(appData).Error;
		appData, err = db.AppData().GetAppData(ctx)
		if err != nil {
			if errors.Is(err, serverErrors.ErrDBRecordNotFound) {
				return nil, serverErrors.NewServerErr(fiber.StatusNotFound, "App data not found")
			}

			logger.Error("Error fetching app data from database", zap.Error(err))
			return nil, serverErrors.NewServerErr(fiber.StatusInternalServerError, "Error fetching app data")
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
