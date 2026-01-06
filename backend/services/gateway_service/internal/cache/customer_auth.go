package cache

import (
	"context"
	"errors"
	"strconv"
	"time"
	"vuecom/gateway/internal/constants"
	"vuecom/gateway/internal/types"
	serverErrors "vuecom/shared/errors/server"
	dbModels "vuecom/shared/models/db"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

func GetCustomerSession(token string, api *types.Api, ctx context.Context) (*dbModels.CustomerSession, error) {
	db := api.Deps.DB
	cache := api.Deps.Redis
	logger := api.Deps.Logger

	cus_session := &dbModels.CustomerSession{}

	err := cache.HGetAll(ctx, constants.CUST_SESS+token).Scan(cus_session)

	if err != nil {
		if !errors.Is(err, redis.Nil) {
			logger.Error("failed to get customer session from cache", zap.Error(err))
			return nil, serverErrors.NewServerErr(fiber.StatusInternalServerError, "Something went wrong while getting your session data. Please try again later.")
		}

		logger.Info("customer session not found in cache, fetching from db")
		cus_session, err = db.Customers().GetSessionByToken(token, ctx)
		if err != nil {
			if errors.Is(err, types.ErrDbNil) {
				logger.Error("customer session not found in db", zap.Error(err))
				return nil, serverErrors.NewServerErr(fiber.StatusUnauthorized, "User Session not found. Consider logging in again")
			}

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

func GetCustomerById(api *types.Api, id int, ctx context.Context) (*dbModels.Customer, error) {
	db := api.Deps.DB
	cache := api.Deps.Redis
	logger := api.Deps.Logger
	customer := &dbModels.Customer{}

	// Try to get from cache first
	err := cache.HGetAll(ctx, constants.CUST_KEY+strconv.Itoa(id)).Scan(customer)
	if err != nil {
		if !errors.Is(err, redis.Nil) {
			logger.Error("failed to get customer from cache", zap.Error(err))
			return nil, serverErrors.NewServerErr(fiber.StatusInternalServerError, "Something went wrong while getting your session data. Please try again later.")
		}

		logger.Info("customer not found in cache, fetching from db")
		customer, err = db.Customers().GetUserById(id, ctx)

		if err != nil {
			if errors.Is(err, types.ErrDbNil) {
				logger.Error("customer not found in db", zap.Error(err))
				return nil, serverErrors.NewServerErr(fiber.StatusUnauthorized, "Customer not found. Consider logging in again")
			}

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
