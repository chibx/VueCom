package main

import (
	_api "vuecom/api"

	"github.com/gofiber/fiber/v2"
)

// Potentially Long Function | Just stack all the routes in here
func loadApis(app fiber.Router, handler *_api.Api) {

	/* /api handlers */
	api := app.Group("/api")
	api.Get("/products", handler.GetProducts)
	api.Post("/app/create")

	// Normal App Handlers
	app.Use("/:admin/*", handler.ValidateSlug)
}
