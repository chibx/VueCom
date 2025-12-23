package gorm_pg

import (
	"context"
	"vuecom/gateway/internal/types"
	dbModels "vuecom/shared/models/db"

	"gorm.io/gorm"
)

type appdataRepository struct {
	db *gorm.DB
}

func (a *appdataRepository) GetAppData(ctx context.Context) (*dbModels.AppData, error) {
	return nil, types.ErrDbUnimplemented
}
