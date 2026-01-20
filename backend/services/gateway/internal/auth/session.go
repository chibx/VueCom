package auth

import (
	"context"
	"crypto/rand"
	"vuecom/gateway/internal/constants"
	"vuecom/gateway/internal/types"
	userModels "vuecom/shared/models/db/users"
)

func GenerateSessionToken() (string, error) {
	// Generate a random session token
	bytes := make([]byte, 32) // e.g., 32 for ~256 bits entropy
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return constants.BU_KEY + string(bytes), nil
}

func GenerateCustomerSessionToken() (string, error) {
	// Generate a random session token
	bytes := make([]byte, 32) // e.g., 32 for ~256 bits entropy
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return constants.CUST_KEY + string(bytes), nil
}

func DeleteBackendSession(ctx context.Context, api *types.Api, session *userModels.BackendSession) error {
	// Delete the session from the database and cache
	db := api.Deps.DB
	cache := api.Deps.Redis
	if err := cache.Unlink(ctx, constants.BU_KEY+session.Token).Err(); err != nil {
		return err
	}

	// err := db.Model(backendSession).Where("token = ?", token).Delete(backendSession).Error
	err := db.BackendUsers().DeleteSession(session, ctx)
	if err != nil {
		return err
	}

	return nil
}

func DeleteCustomerSession(ctx context.Context, api *types.Api, session *userModels.CustomerSession) error {
	// Delete the session from the database and cache
	db := api.Deps.DB
	cache := api.Deps.Redis
	if err := cache.Unlink(ctx, constants.CUST_KEY+session.Token).Err(); err != nil {
		return err
	}

	// err := db.Model(customerSession).Where("token = ?", token).Delete(customerSession).Error
	err := db.Customers().DeleteSession(session, ctx)
	if err != nil {
		return err
	}

	return nil
}
