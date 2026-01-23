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

func GetCustomerSession(token string, api *types.Api, ctx context.Context) (*userModels.CustomerSession, error) {
	db := api.Deps.DB
	cache := api.Deps.Redis
	logger := api.Deps.Logger

	cus_session := &userModels.CustomerSession{}

	err := cache.HGetAll(ctx, constants.CUST_SESS+token).Scan(cus_session)
	notExist := cus_session.UserID == 0

	if err != nil || notExist {
		if notExist {
			logger.Info("customer session not found in cache, fetching from db")
			// return nil, serverErrors.NewServerErr(fiber.StatusInternalServerError, "Something went wrong while getting your session data. Please try again later.")
		} else {
			logger.Error("failed to get customer session from cache", zap.Error(err))
		}

		cus_session, err = db.Customers().GetSessionByToken(ctx, token)

		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Error("customer session not found in db", zap.Error(err))
			return nil, serverErrors.NewServerErr(fiber.StatusUnauthorized, "User Session not found. Consider logging in again")
		}

		if err != nil {
			logger.Error("failed to get customer session from db", zap.Error(err))
			return nil, serverErrors.NewServerErr(fiber.StatusInternalServerError, "Something went wrong while fetching your session data. Please try again later.")
		}

		go func() {
			logger.Info("caching customer session")
			_, err := cache.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
				pipe.HSet(ctx, constants.CUST_SESS+token, cus_session)
				pipe.Expire(ctx, constants.CUST_SESS+token, 5*time.Minute) // Global expiry on the key.
				return nil
			})

			if err != nil {
				logger.Error("failed to cache customer session", zap.Error(err))
			}
		}()
	}

	return cus_session, nil
}

func GetCustomerById(api *types.Api, id int, ctx context.Context) (*userModels.Customer, error) {
	db := api.Deps.DB
	cache := api.Deps.Redis
	logger := api.Deps.Logger
	customer := &userModels.Customer{}

	// Try to get from cache first
	err := cache.HGetAll(ctx, constants.CUST_KEY+strconv.Itoa(id)).Scan(customer)
	notExist := customer.ID == 0

	if err != nil || notExist {
		if notExist {
			logger.Info("customer not found in cache, fetching from db")
			// return nil, serverErrors.NewServerErr(fiber.StatusInternalServerError, "Something went wrong while getting your session data. Please try again later.")
		} else {
			logger.Error("failed to get customer from cache", zap.Error(err))
		}

		customer, err = db.Customers().GetUserById(ctx, id)

		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Error("customer not found in db", zap.Error(err))
			return nil, serverErrors.NewServerErr(fiber.StatusUnauthorized, "Customer not found. Consider logging in again")
		}

		if err != nil {
			logger.Error("failed to get customer from db", zap.Error(err))
			return nil, serverErrors.NewServerErr(fiber.StatusInternalServerError, "Something went wrong while fetching your session data. Please try again later.")
		}

		go func() {
			logger.Info("caching customer")
			_, err := cache.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
				pipe.HSet(ctx, constants.CUST_KEY+strconv.Itoa(id), customer)
				pipe.Expire(ctx, constants.CUST_KEY+strconv.Itoa(id), 5*time.Minute) // Global expiry on the key.
				return nil
			})

			if err != nil {
				logger.Error("failed to cache customer", zap.Error(err))
			}
		}()
	}

	return customer, nil
}

func TouchCustomerSession(api *types.Api, token string, ctx context.Context) {
	rdb := api.Deps.Redis
	logger := api.Deps.Logger

	// ttl := rdb.TTL(ctx, token)
	// if ttl.Err() != nil {
	// 	logger.Error("failed to get backend user session ttl", zap.Error(ttl.Err()))
	// 	return
	// }
	// ttlSeconds := int64(ttl.Val() / time.Second)
	// logger.Info("backend user session ttl", zap.Int64("ttl", ttlSeconds))

	_, err := rdb.Expire(ctx, constants.CUST_SESS+token, constants.BackendSessionTimeout).Result()
	if err != nil {
		logger.Error("failed to expire customer session", zap.Error(err))
	}
}
