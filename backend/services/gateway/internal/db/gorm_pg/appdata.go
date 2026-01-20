package gorm_pg

import (
	"context"

	appModels "github.com/chibx/vuecom/backend/shared/models/db/appdata"
	userModels "github.com/chibx/vuecom/backend/shared/models/db/users"

	"gorm.io/gorm"
)

type appdataRepository struct {
	db *gorm.DB
}

func (ar *appdataRepository) CreateAppData(appData *appModels.AppData, ctx context.Context) error {
	return ar.db.WithContext(ctx).Create(appData).Error
}

func (ar *appdataRepository) CountOwner(ctx context.Context) (int64, error) {
	count := int64(0)
	err := ar.db.WithContext(ctx).Model(&userModels.BackendUser{}).Where("role = 'owner'").Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
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
