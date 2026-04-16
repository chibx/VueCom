package cache

import (
	"context"
	"errors"
	"strconv"
	"time"

	userModels "github.com/chibx/vuecom/backend/shared/models/db/users"

	"github.com/chibx/vuecom/backend/services/gateway/internal/cache/keys"
	"github.com/chibx/vuecom/backend/services/gateway/internal/global"
	"github.com/chibx/vuecom/backend/services/gateway/internal/types"

	serverErrors "github.com/chibx/vuecom/backend/shared/errors/server"
	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

func GetBackendUserById(api *types.Api, id uint32, ctx context.Context) (*userModels.BackendUser, error) {
	db := api.Deps.DB
	cache := api.Deps.Redis
	logger := global.Logger
	backendUser := &userModels.BackendUser{}

	// Try to get from cache first
	err := cache.HGetAll(ctx, keys.BackendUserKey(id)).Scan(backendUser)
	notExist := backendUser.ID == 0

	if err != nil || notExist {
		if notExist {
			logger.Info("backend user not found in cache, fetching from db")
		} else {
			logger.Error("failed to get backend user from cache", zap.Error(err))
		}

		backendUser, err = db.BackendUsers().GetUserById(ctx, id)

		if errors.Is(err, serverErrors.ErrDBRecordNotFound) {
			logger.Error("backend user" + strconv.Itoa(int(id)) + "not found in db")
			return nil, serverErrors.NewServerErr(fiber.StatusUnauthorized, "User ID "+strconv.Itoa(int(id))+" not found. Consider logging in again")
		}

		if err != nil {
			// err can't be serverErrors.ErrDBRecordNotFound so it's fine here
			logger.Error("failed to get backend user from db", zap.Error(err))
			return nil, serverErrors.NewServerErr(fiber.StatusInternalServerError, "Something went wrong while fetching your session data. Please try again later.")
		}

		go func() {
			logger.Info("caching backend user")
			_, err := cache.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
				var err error
				err = pipe.HSet(ctx, keys.BackendUserKey(id), backendUser).Err()
				if err != nil {
					return err
				}
				err = pipe.Expire(ctx, keys.BackendUserKey(id), 5*time.Minute).Err() // Global expiry on the key.
				return err
			})

			if err != nil {
				logger.Error("failed to cache backend user", zap.Error(err))
			}
		}()
	}

	return backendUser, nil
}
