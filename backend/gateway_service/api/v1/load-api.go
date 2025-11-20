package v1

import (
	"github.com/gofiber/fiber/v2"
)

// Potentially Long Function | Just stack all the routes in here
func (v1_api *Api) LoadApis(app fiber.Router) {

	/* /api handlers */
	api := app.Group("/api")
	api.Post("/product", v1_api.CreateProduct)
	api.Get("/product/:id", v1_api.GetProduct)
	api.Post("/app/admin-exist")

	// Normal App Handlers
	app.Use("/:admin/*", v1_api.ValidateSlug)
}
