package products

import (
	"vuecom/gateway/internal/types"

	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(app fiber.Router, api *types.Api) {
	app.Post("/product", func(ctx *fiber.Ctx) error {
		return CreateProduct(ctx, api)
	})
	app.Get("/product/:id", func(ctx *fiber.Ctx) error {
		return GetProduct(ctx, api)
	})
}
