package cache

import (
	"errors"
	"vuecom/gateway/internal/cache/keys"
	"vuecom/gateway/internal/v1/types"
	dbModels "vuecom/shared/models/db"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func GetAppData(ctx *fiber.Ctx, api *types.Api) (*dbModels.AppData, error) {
	db := api.Deps.DB
	cache := api.Deps.Redis
	context := ctx.Context()
	appData := new(dbModels.AppData)

	err := cache.Get(context, keys.APP_DATA_KEY).Scan(&appData)
	if err != nil {
		if !errors.Is(err, redis.Nil) {
			return nil, fiber.NewError(fiber.StatusInternalServerError, "Error fetching app data")
		}

		// Key Not Found
		if err := db.WithContext(context).Limit(1).First(appData).Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fiber.NewError(fiber.StatusInternalServerError, "Error fetching app data")
		}

		if appData == nil {
			return nil, nil
		}

		if err := cache.Set(context, keys.APP_DATA_KEY, appData, 0).Err(); err != nil {
			// TODO: Log cache issue
			return nil, fiber.NewError(fiber.StatusInternalServerError, "Error fetching app data")
		}

	}

	return appData, nil
}
