package cache

import (
	"context"
	"errors"
	"strconv"
	"time"

	serverErrors "github.com/chibx/vuecom/backend/shared/errors/server"
	userModels "github.com/chibx/vuecom/backend/shared/models/db/users"
	"gorm.io/gorm"

	"github.com/chibx/vuecom/backend/services/gateway/internal/constants"
	"github.com/chibx/vuecom/backend/services/gateway/internal/types"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

func GetBackendUserSession(token string, api *types.Api, ctx context.Context) (*userModels.BackendSession, error) {
	db := api.Deps.DB
	cache := api.Deps.Redis
	logger := api.Deps.Logger

	var backend_session *userModels.BackendSession

	err := cache.HGetAll(ctx, constants.BU_SESS+token).Scan(&backend_session)
	notExist := backend_session.UserId == 0

	if err != nil || notExist {
		if notExist {
			logger.Info("backend user session not found in cache")
		} else {
			logger.Error("failed to get backend user session from cache, fetching from db", zap.Error(err))
		}

		backend_session, err = db.BackendUsers().GetSessionByToken(ctx, token)

		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Error("backend user session not found in db", zap.Error(err))
			return nil, serverErrors.NewServerErr(fiber.StatusUnauthorized, "User Session not found. Consider logging in again")
		}

		if err != nil {
			// err can't be ErrRecordNotFound so it's fine here
			logger.Error("failed to get backend user session from db", zap.Error(err))
			return nil, serverErrors.NewServerErr(fiber.StatusInternalServerError, "Something went wrong while fetching your session data. Please try again later.")
		}

		go func() {
			logger.Info("caching backend user session")
			_, err := cache.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
				pipe.HSet(ctx, constants.BU_SESS+token, backend_session)
				pipe.Expire(ctx, constants.BU_SESS+token, 5*time.Minute) // Global expiry on the key.
				return nil
			})

			if err != nil {
				logger.Error("failed to cache backend user session", zap.Error(err))
			}
		}()
	}

	return backend_session, nil
}

func GetBackendUserById(api *types.Api, id int, ctx context.Context) (*userModels.BackendUser, error) {
	db := api.Deps.DB
	cache := api.Deps.Redis
	logger := api.Deps.Logger
	backendUser := &userModels.BackendUser{}

	// Try to get from cache first
	err := cache.HGetAll(ctx, constants.BU_KEY+strconv.Itoa(id)).Scan(backendUser)
	notExist := backendUser.ID == 0

	if err != nil || notExist {
		if notExist {
			logger.Info("backend user not found in cache, fetching from db")
		} else {
			logger.Error("failed to get backend user from cache", zap.Error(err))
		}

		backendUser, err = db.BackendUsers().GetUserById(ctx, id)

		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Error("backend user" + strconv.Itoa(id) + "not found in db")
			return nil, serverErrors.NewServerErr(fiber.StatusUnauthorized, "User ID "+strconv.Itoa(id)+" not found. Consider logging in again")
		}

		if err != nil {
			// err can't be ErrRecordNotFound so it's fine here
			logger.Error("failed to get backend user from db", zap.Error(err))
			return nil, serverErrors.NewServerErr(fiber.StatusInternalServerError, "Something went wrong while fetching your session data. Please try again later.")
		}

		go func() {
			logger.Info("caching backend user")
			_, err := cache.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
				pipe.HSet(ctx, constants.BU_KEY+strconv.Itoa(id), backendUser)
				pipe.Expire(ctx, constants.BU_KEY+strconv.Itoa(id), 5*time.Minute) // Global expiry on the key.
				return nil
			})

			if err != nil {
				logger.Error("failed to cache backend user", zap.Error(err))
			}
		}()
	}

	return backendUser, nil
}

func TouchBackendSession(api *types.Api, token string, ctx context.Context) {
	rdb := api.Deps.Redis
	logger := api.Deps.Logger

	// ttl := rdb.TTL(ctx, token)
	// if ttl.Err() != nil {
	// 	logger.Error("failed to get backend user session ttl", zap.Error(ttl.Err()))
	// 	return
	// }
	// ttlSeconds := int64(ttl.Val() / time.Second)
	// logger.Info("backend user session ttl", zap.Int64("ttl", ttlSeconds))

	_, err := rdb.Expire(ctx, constants.BU_SESS+token, constants.BackendSessionTimeout).Result()
	if err != nil {
		logger.Error("failed to expire backend user session", zap.Error(err))
	}
}
