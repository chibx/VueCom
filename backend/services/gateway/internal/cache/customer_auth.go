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
