package main

import (
	_api "vuecom/api"

	"github.com/gofiber/fiber/v2"
)

// Potentially Long Function | Just stack all the routes in here
func LoadApis(app fiber.Router) {
	handler := _api.Api{}

	/* /api handlers */
	api := app.Group("/api")
	api.Get("/:name", handler.ApiHandler)
	api.Get("/products", handler.GetProducts)

	app.Use("/:admin/*", handler.ValidateSlug)
}
