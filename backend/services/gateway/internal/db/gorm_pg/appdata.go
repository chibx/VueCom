package gorm_pg

import (
	"context"

	appModels "github.com/chibx/vuecom/backend/shared/models/db/appdata"

	"gorm.io/gorm"
)

type appdataRepository struct {
	db *gorm.DB
}

func (ar *appdataRepository) CreateAppData(ctx context.Context, appData *appModels.AppData) error {
	return ar.db.WithContext(ctx).Create(appData).Error
}

// err := api.Deps.DB.First(appData).Error

func (ar *appdataRepository) GetAppData(ctx context.Context) (*appModels.AppData, error) {
	appData := &appModels.AppData{}
	err := ar.db.WithContext(ctx).First(appData).Error
	if err != nil {
		return nil, err
	}
	return appData, nil
}
