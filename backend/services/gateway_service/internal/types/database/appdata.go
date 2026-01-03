package database

import (
	"context"
	dbModels "vuecom/shared/models/db"
)

type AppDataRepository interface {
	CreateAppData(appData *dbModels.AppData, ctx context.Context) error
	GetAppData(ctx context.Context) (*dbModels.AppData, error)
	// api.Deps.DB.Model(&dbModels.BackendUser{}).Where("role = 'owner'").Count(&count)
	CountOwner(ctx context.Context) (int64, error)
}
