package orders

import (
	"github.com/chibx/vuecom/backend/services/gateway/internal/types"

	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(app fiber.Router, api *types.Api) {
	app.Post("/order", func(ctx *fiber.Ctx) error {
		return CreateOrder(ctx, api)
	})
	app.Get("/order/:id", func(ctx *fiber.Ctx) error {
		return GetOrder(ctx, api)
	})
	app.Delete("/order/:id", func(ctx *fiber.Ctx) error {
		return DeleteOrder(ctx, api)
	})
	app.Delete("/orders", func(ctx *fiber.Ctx) error {
		return DeleteOrders(ctx, api)
	})
	app.Put("/order/:id", func(ctx *fiber.Ctx) error {
		return UpdateOrder(ctx, api)
	})
	app.Get("/orders", func(ctx *fiber.Ctx) error {
		return ListOrders(ctx, api)
	})
}
