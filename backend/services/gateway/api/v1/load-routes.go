package v1

import (
	"github.com/chibx/vuecom/backend/services/gateway/api/v1/handlers/backend"
	"github.com/chibx/vuecom/backend/services/gateway/api/v1/handlers/customer"
	"github.com/chibx/vuecom/backend/services/gateway/api/v1/middlewares"

	// "github.com/chibx/vuecom/backend/services/gateway/internal/constants"
	"github.com/chibx/vuecom/backend/services/gateway/internal/constants"
	"github.com/chibx/vuecom/backend/services/gateway/internal/types"

	"github.com/gofiber/fiber/v2"
)

// Potentially Long Function | Just stack all the routes in here
func LoadRoutes(app fiber.Router, api *types.Api) {
	// app.Use(middlewares.AuthMiddleware(api), middlewares.ServeIndex(api), middlewares.ServeAssets())
	app.Get("/api/health", func(ctx *fiber.Ctx) error {
		return ctx.Status(200).SendString("OK")
	})

	app.Use(middlewares.ServeAssets(), middlewares.AuthMiddleware(api))

	backend.LoadRoutes(app, api)
	customer.LoadRoutes(app, api)

	app.Use(middlewares.ServeIndex(api))

	// app.Static("*", "./"+constants.PublicFolder, fiber.Static{})
	app.Get("*", func(ctx *fiber.Ctx) error {
		return ctx.SendFile("./" + constants.PublicFolder + "/index.html")
	})
}
