package v1

import (
	"vuecom/gateway/api/v1/handlers"
	adminHandler "vuecom/gateway/api/v1/handlers/admin"
	orderHandler "vuecom/gateway/api/v1/handlers/orders"
	productHandler "vuecom/gateway/api/v1/handlers/products"
	"vuecom/gateway/internal/v1/types"

	"github.com/gofiber/fiber/v2"
)

// Potentially Long Function | Just stack all the routes in here
func LoadRoutes(app fiber.Router, api *types.Api) {
	/**
	  * 	func(ctx *fiber.Ctx) error {
	  		return
	    	}
	  *
	*/

	/* /v1 handlers */
	v1 := app.Group("/api/v1")
	productHandler.RegisterRoutes(v1, api)
	adminHandler.RegisterRoutes(v1, api)
	orderHandler.RegisterRoutes(v1, api)

	// Normal App Handlers
	app.Use("/:admin/*", func(ctx *fiber.Ctx) error {
		return handlers.ValidateSlug(ctx, api)
	})
}
