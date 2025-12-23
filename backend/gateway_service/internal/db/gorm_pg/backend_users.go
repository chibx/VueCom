package gorm_pg

import (
	"context"
	"errors"
	"vuecom/gateway/internal/types"
	dbModels "vuecom/shared/models/db"

	"gorm.io/gorm"
)

type backendUserRepository struct {
	db *gorm.DB
}

func (b *backendUserRepository) GetBackendUserById(id int, ctx context.Context) (*dbModels.BackendUser, error) {
	backendUser := &dbModels.BackendUser{}
	result := b.db.WithContext(ctx).First(backendUser, "id = ?", id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, types.ErrDbNil
		}

		return nil, result.Error
	}

	return backendUser, nil
}

func (b *backendUserRepository) GetBackendUserByApiKey(apiKey string, ctx context.Context) (*dbModels.BackendUser, error) {
	return nil, types.ErrDbUnimplemented
}
