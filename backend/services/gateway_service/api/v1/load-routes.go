package v1

import (
	"vuecom/gateway/api/v1/handlers/backend"
	"vuecom/gateway/api/v1/handlers/customer"
	"vuecom/gateway/api/v1/middlewares"
	"vuecom/gateway/internal/constants"
	"vuecom/gateway/internal/types"

	"github.com/gofiber/fiber/v2"
)

// Potentially Long Function | Just stack all the routes in here
func LoadRoutes(app fiber.Router, api *types.Api) {
	// app.Use(middlewares.AuthMiddleware(api), middlewares.ServeIndex(api), middlewares.ServeAssets())
	app.Get("/api/health", func(ctx *fiber.Ctx) error {
		return ctx.Status(200).SendString("OK")
	})

	app.Use(middlewares.AuthMiddleware(api), middlewares.ServeIndex(api))

	backend.LoadRoutes(app, api)
	customer.LoadRoutes(app, api)

	app.Static("*", "./"+constants.PublicFolder)
	app.Get("/*", func(ctx *fiber.Ctx) error {
		return ctx.SendFile("./" + constants.PublicFolder + "/index.html")
	})
}
