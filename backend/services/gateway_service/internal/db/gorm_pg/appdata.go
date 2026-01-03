package gorm_pg

import (
	"context"
	dbModels "vuecom/shared/models/db"

	"gorm.io/gorm"
)

type appdataRepository struct {
	db *gorm.DB
}

func (ar *appdataRepository) CreateAppData(appData *dbModels.AppData, ctx context.Context) error {
	return ar.db.WithContext(ctx).Create(appData).Error
}

func (ar *appdataRepository) CountOwner(ctx context.Context) (int64, error) {
	count := int64(0)
	err := ar.db.WithContext(ctx).Model(&dbModels.BackendUser{}).Where("role = 'owner'").Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

// err := api.Deps.DB.First(appData).Error

func (ar *appdataRepository) GetAppData(ctx context.Context) (*dbModels.AppData, error) {
	appData := &dbModels.AppData{}
	err := ar.db.WithContext(ctx).First(appData).Error
	if err != nil {
		return nil, err
	}
	return appData, nil
}
