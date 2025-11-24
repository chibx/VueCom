package v1

import (
	"vuecom/gateway/api/v1/handlers"
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
	v1.Post("/product", func(ctx *fiber.Ctx) error {
		return handlers.CreateProduct(ctx, api)
	})
	v1.Get("/product/:id", func(ctx *fiber.Ctx) error {
		return handlers.GetProduct(ctx, api)
	})
	v1.Get("/admin-exist", func(ctx *fiber.Ctx) error {
		exists, err := handlers.DoesOwnerExist(ctx, api)
		if err != nil {
			return fiber.ErrInternalServerError
		}
		return ctx.JSON(fiber.Map{
			"exists": exists,
		})
	})

	// Normal App Handlers
	app.Use("/:admin/*", func(ctx *fiber.Ctx) error {
		return handlers.ValidateSlug(ctx, api)
	})
}
