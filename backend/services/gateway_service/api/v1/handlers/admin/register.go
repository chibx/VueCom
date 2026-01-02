package admin

import (
	"vuecom/gateway/internal/types"

	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(app fiber.Router, api *types.Api) {
	app.Post("/initialize-app", func(ctx *fiber.Ctx) error {
		return InitializeApp(ctx, api)
	})
	app.Post("/register-owner", func(ctx *fiber.Ctx) error {
		return RegisterOwner(ctx, api)
	})
	app.Get("/admin-exist", func(ctx *fiber.Ctx) error {
		exists, err := DoesOwnerExist(ctx, api)
		if err != nil {
			return fiber.ErrInternalServerError
		}
		return ctx.JSON(fiber.Map{
			"exists": exists,
		})
	})
}
