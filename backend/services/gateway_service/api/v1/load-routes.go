package v1

import (
	adminHandler "vuecom/gateway/api/v1/handlers/admin"
	orderHandler "vuecom/gateway/api/v1/handlers/orders"
	productHandler "vuecom/gateway/api/v1/handlers/products"
	"vuecom/gateway/api/v1/middlewares"
	"vuecom/gateway/internal/types"

	"github.com/gofiber/fiber/v2"
)

// Potentially Long Function | Just stack all the routes in here
func LoadRoutes(app fiber.Router, api *types.Api) {
	app.Use(middlewares.AuthMiddleware(api), middlewares.ServeAssets())

	/* /v1 handlers */
	v1 := app.Group("/api/v1")
	productHandler.RegisterRoutes(v1, api)
	adminHandler.RegisterRoutes(v1, api)
	orderHandler.RegisterRoutes(v1, api)

	// Normal App Handlers

	// app.Static("*", "./dist", fiber.Static{
	// 	Next: func(ctx *fiber.Ctx) bool {
	// 		logger := api.Deps.Logger
	// 		routeParts := utils.ExtractRouteParts(ctx.Path())

	// 		if len(routeParts) > 1 && routeParts[1] == api.AdminSlug {
	// 			logger.Info("Admin route detected", zap.String("route", ctx.Path()))
	// 			return true
	// 		}

	// 		return false
	// 	},
	// })
}
