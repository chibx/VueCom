package auth

import (
	"context"
	"crypto/rand"
	"vuecom/gateway/internal/v1/types"
	dbModels "vuecom/shared/models/db"
)

func GenerateSessionToken() (string, error) {
	// Generate a random session token
	bytes := make([]byte, 32) // e.g., 32 for ~256 bits entropy
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return "session:" + string(bytes), nil
}

func DeleteBackendSession(ctx context.Context, api *types.Api, token string) error {
	// Delete the session from the database and cache
	db := api.Deps.DB
	cache := api.Deps.Redis
	backendSession := &dbModels.BackendSession{}
	if err := cache.Unlink(ctx, token).Err(); err != nil {
		return err
	}

	err := db.Model(backendSession).Where("token = ?", token).Delete(backendSession).Error
	if err != nil {
		return err
	}

	return nil
}
