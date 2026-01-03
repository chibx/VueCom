package bgtasks

import (
	"context"
	"time"
	"vuecom/gateway/internal/cache"
	"vuecom/gateway/internal/types"
)

// This helps to refresh the local copy of the slug from the data fetching function
func RefreshAdminSlug(api *types.Api) {
	ticker := time.NewTicker(time.Minute)
	ctx := context.Background()
	go func() {
		for range ticker.C {
			appData, _ := cache.GetAppData(ctx, api)

			if appData != nil {
				api.AdminSlug = appData.AdminRoute
			}
		}
	}()
}
