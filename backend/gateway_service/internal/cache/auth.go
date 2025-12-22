package cache

import (
	"context"
	"errors"
	"vuecom/gateway/internal/v1/types"
	userErrors "vuecom/shared/errors/users"
	dbModels "vuecom/shared/models/db"

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
			return nil, userErrors.NewTokenErr(500, "Something went wrong while getting your session data. Please try again later.")
		}

		result := db.WithContext(context).First(cus_session, "token = ?", token)
		if result.Error != nil {
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				return nil, userErrors.NewTokenErr(401, "User Session not found. Consider logging in again")
			}

			return nil, userErrors.NewTokenErr(500, "Something went wrong while fetching your session data. Please try again later.")
		}

		go func() {
			err = cache.HSet(context, "c_sess:"+token, cus_session).Err()
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
			return nil, userErrors.NewTokenErr(500, "Something went wrong while getting your session data. Please try again later.")
		}

		result := db.WithContext(context).First(backend_session, "token = ?", token)
		if result.Error != nil {
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				return nil, userErrors.NewTokenErr(401, "User Session not found. Consider logging in again")
			}

			return nil, userErrors.NewTokenErr(500, "Something went wrong while fetching your session data. Please try again later.")
		}

		go func() {
			err = cache.HSet(context, "b_sess:"+token, backend_session).Err()
			if err != nil {
				// Log the error but don't fail the request
				// The user session is still returned from the database
				// TODO: Add logging here
			}
		}()
	}

	return backend_session, nil
}
