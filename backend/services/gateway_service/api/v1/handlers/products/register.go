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
	app.Put("/product/:id", func(ctx *fiber.Ctx) error {
		return UpdateProduct(ctx)
	})
	app.Delete("/product/:id", func(ctx *fiber.Ctx) error {
		return DeleteProduct(ctx)
	})
	app.Delete("/products", func(ctx *fiber.Ctx) error {
		return DeleteProducts(ctx, api)
	})
}
