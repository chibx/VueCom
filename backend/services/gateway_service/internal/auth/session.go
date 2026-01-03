package auth

import (
	"context"
	"crypto/rand"
	"vuecom/gateway/internal/types"
	dbModels "vuecom/shared/models/db"
)

func GenerateSessionToken() (string, error) {
	// Generate a random session token
	bytes := make([]byte, 32) // e.g., 32 for ~256 bits entropy
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return "b_sess:" + string(bytes), nil
}

func GenerateCustomerSessionToken() (string, error) {
	// Generate a random session token
	bytes := make([]byte, 32) // e.g., 32 for ~256 bits entropy
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return "c_sess:" + string(bytes), nil
}

func DeleteBackendSession(ctx context.Context, api *types.Api, session *dbModels.BackendSession) error {
	// Delete the session from the database and cache
	db := api.Deps.DB
	cache := api.Deps.Redis
	if err := cache.Unlink(ctx, "b_sess:"+session.Token).Err(); err != nil {
		return err
	}

	// err := db.Model(backendSession).Where("token = ?", token).Delete(backendSession).Error
	err := db.BackendUsers().DeleteSession(session, ctx)
	if err != nil {
		return err
	}

	return nil
}

func DeleteCustomerSession(ctx context.Context, api *types.Api, session *dbModels.CustomerSession) error {
	// Delete the session from the database and cache
	db := api.Deps.DB
	cache := api.Deps.Redis
	if err := cache.Unlink(ctx, "c_sess:"+session.Token).Err(); err != nil {
		return err
	}

	// err := db.Model(customerSession).Where("token = ?", token).Delete(customerSession).Error
	err := db.Customers().DeleteSession(session, ctx)
	if err != nil {
		return err
	}

	return nil
}
