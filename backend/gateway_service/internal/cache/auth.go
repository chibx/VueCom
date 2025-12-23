package cache

import (
	"context"
	"errors"
	"strconv"
	"time"
	"vuecom/gateway/internal/types"
	serverErrors "vuecom/shared/errors/server"
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
			return nil, serverErrors.NewTokenErr(fiber.StatusInternalServerError, "Something went wrong while getting your session data. Please try again later.")
		}

		result := db.WithContext(context).First(cus_session, "token = ?", token)
		if result.Error != nil {
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				return nil, serverErrors.NewTokenErr(fiber.StatusUnauthorized, "User Session not found. Consider logging in again")
			}

			return nil, serverErrors.NewTokenErr(fiber.StatusInternalServerError, "Something went wrong while fetching your session data. Please try again later.")
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
			return nil, serverErrors.NewTokenErr(fiber.StatusInternalServerError, "Something went wrong while getting your session data. Please try again later.")
		}

		result := db.WithContext(context).First(backend_session, "token = ?", token)
		if result.Error != nil {
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				return nil, serverErrors.NewTokenErr(fiber.StatusUnauthorized, "User Session not found. Consider logging in again")
			}

			return nil, serverErrors.NewTokenErr(fiber.StatusInternalServerError, "Something went wrong while fetching your session data. Please try again later.")
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

func GetBackendUserById(api *types.Api, id int, context context.Context) (*dbModels.BackendUser, error) {
	db := api.Deps.DB
	cache := api.Deps.Redis
	backendUser := &dbModels.BackendUser{}

	// Try to get from cache first
	err := cache.HGetAll(context, "b_user:"+strconv.Itoa(id)).Scan(backendUser)
	if err != nil {
		if !errors.Is(err, redis.Nil) {
			return nil, serverErrors.NewTokenErr(fiber.StatusInternalServerError, "Something went wrong while getting your session data. Please try again later.")
		}

		result := db.WithContext(context).First(backendUser, "id = ?", id)
		if result.Error != nil {
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				return nil, serverErrors.NewTokenErr(fiber.StatusUnauthorized, "User not found. Consider logging in again")
			}

			return nil, serverErrors.NewTokenErr(fiber.StatusInternalServerError, "Something went wrong while fetching your session data. Please try again later.")
		}

		go func() {
			_, err := cache.TxPipelined(context, func(pipe redis.Pipeliner) error {
				pipe.HSet(context, "b_user:"+strconv.Itoa(id), backendUser)
				pipe.Expire(context, "b_user:"+strconv.Itoa(id), 10*time.Minute) // Global expiry on the key.
				return nil
			})

			// err = cache.HSet(context, "b_user:"+strconv.Itoa(id), backendUser).Err()
			if err != nil {
				// Log the error but don't fail the request
				// The user session is still returned from the database
				// TODO: Add logging here
			}
		}()
	}

	return backendUser, nil
}

func GetCustomerById(api *types.Api, id int, context context.Context) (*dbModels.Customer, error) {
	db := api.Deps.DB
	cache := api.Deps.Redis
	customer := &dbModels.Customer{}

	// Try to get from cache first
	err := cache.HGetAll(context, "cust:"+strconv.Itoa(id)).Scan(customer)
	if err != nil {
		if !errors.Is(err, redis.Nil) {
			return nil, serverErrors.NewTokenErr(fiber.StatusInternalServerError, "Something went wrong while getting your session data. Please try again later.")
		}

		result := db.WithContext(context).First(customer, "id = ?", id)
		if result.Error != nil {
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				return nil, serverErrors.NewTokenErr(fiber.StatusUnauthorized, "Customer not found. Consider logging in again")
			}

			return nil, serverErrors.NewTokenErr(fiber.StatusInternalServerError, "Something went wrong while fetching your session data. Please try again later.")
		}

		go func() {
			_, err := cache.TxPipelined(context, func(pipe redis.Pipeliner) error {
				pipe.HSet(context, "cust:"+strconv.Itoa(id), customer)
				pipe.Expire(context, "cust:"+strconv.Itoa(id), 10*time.Minute) // Global expiry on the key.
				return nil
			})

			// err = cache.HSet(context, "cust:"+strconv.Itoa(id), customer).Err()
			if err != nil {
				// Log the error but don't fail the request
				// The user session is still returned from the database
				// TODO: Add logging here
			}
		}()
	}

	return customer, nil
}
