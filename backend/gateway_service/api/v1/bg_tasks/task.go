package bgtasks

import (
	"context"
	"time"
	"vuecom/gateway/internal/cache"
	"vuecom/gateway/internal/v1/types"
)

// This helps to refresh the local copy of the slug from the data fetching function
func RefreshAdminSlug(api *types.Api) {
	ticker := time.NewTicker(time.Minute)

	go func() {
		for range ticker.C {
			appData, _ := cache.GetAppData(context.Background(), api)

			if appData != nil {
				api.AdminSlug = appData.AdminRoute
			}
		}
	}()
}
