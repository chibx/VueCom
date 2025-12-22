package cache

import (
	"context"
	"errors"
	"time"
	"vuecom/gateway/internal/v1/types"
	userErrors "vuecom/shared/errors/users"
	dbModels "vuecom/shared/models/db"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func GetCustomerSession(token string, api *types.Api, context context.Context) (*dbModels.CustomerSession, error) {
	db := api.Deps.DB
	cache := api.Deps.Redis

	cus_session := &dbModels.CustomerSession{}

	err := cache.HGetAll(context, "c_sess:"+token).Scan(cus_session)

	if err != nil {
		if !errors.Is(err, redis.Nil) {
			return nil, userErrors.NewTokenErr(fiber.StatusInternalServerError, "Something went wrong while getting your session data. Please try again later.")
		}

		result := db.WithContext(context).First(cus_session, "token = ?", token)
		if result.Error != nil {
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				return nil, userErrors.NewTokenErr(fiber.StatusUnauthorized, "User Session not found. Consider logging in again")
			}

			return nil, userErrors.NewTokenErr(fiber.StatusInternalServerError, "Something went wrong while fetching your session data. Please try again later.")
		}

		go func() {
			_, err := cache.TxPipelined(context, func(pipe redis.Pipeliner) error {
				pipe.HSet(context, "c_sess:"+token, cus_session)
				pipe.Expire(context, "c_sess:"+token, 5*time.Minute) // Global expiry on the key.
				return nil
			})

			// err = cache.HSet(context, "c_sess:"+token, cus_session).Err()
			if err != nil {
				// Log the error but don't fail the request
				// The user session is still returned from the database
				// TODO: Add logging here
			}
		}()
	}

	return cus_session, nil
}

func GetBackendUserSession(token string, api *types.Api, context context.Context) (*dbModels.BackendSession, error) {
	db := api.Deps.DB
	cache := api.Deps.Redis

	backend_session := &dbModels.BackendSession{}

	err := cache.HGetAll(context, "b_sess:"+token).Scan(backend_session)

	if err != nil {
		if !errors.Is(err, redis.Nil) {
			return nil, userErrors.NewTokenErr(fiber.StatusInternalServerError, "Something went wrong while getting your session data. Please try again later.")
		}

		result := db.WithContext(context).First(backend_session, "token = ?", token)
		if result.Error != nil {
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				return nil, userErrors.NewTokenErr(fiber.StatusUnauthorized, "User Session not found. Consider logging in again")
			}

			return nil, userErrors.NewTokenErr(fiber.StatusInternalServerError, "Something went wrong while fetching your session data. Please try again later.")
		}

		go func() {
			_, err := cache.TxPipelined(context, func(pipe redis.Pipeliner) error {
				pipe.HSet(context, "b_sess:"+token, backend_session)
				pipe.Expire(context, "b_sess:"+token, 5*time.Minute) // Global expiry on the key.
				return nil
			})

			// err = cache.HSet(context, "b_sess:"+token, backend_session).Err()
			if err != nil {
				// Log the error but don't fail the request
				// The user session is still returned from the database
				// TODO: Add logging here
			}
		}()
	}

	return backend_session, nil
}
