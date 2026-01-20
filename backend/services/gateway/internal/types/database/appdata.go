package database

import (
	"context"

	appModels "github.com/chibx/vuecom/backend/shared/models/db/appdata"
)

type AppDataRepository interface {
	CreateAppData(appData *appModels.AppData, ctx context.Context) error
	GetAppData(ctx context.Context) (*appModels.AppData, error)
	// api.Deps.DB.Model(&dbModels.BackendUser{}).Where("role = 'owner'").Count(&count)
	CountOwner(ctx context.Context) (int64, error)
}
