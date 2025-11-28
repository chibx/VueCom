package orders

import (
	"vuecom/gateway/internal/v1/types"

	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(app fiber.Router, api *types.Api) {
	app.Post("/order", func(ctx *fiber.Ctx) error {
		return CreateOrder(ctx, api)
	})
	app.Get("/order/:id", func(ctx *fiber.Ctx) error {
		return GetOrder(ctx, api)
	})
}
