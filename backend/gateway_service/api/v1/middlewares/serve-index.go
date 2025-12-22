package middlewares

import (
	"fmt"
	"vuecom/gateway/internal/utils"
	"vuecom/gateway/internal/v1/types"

	"github.com/gofiber/fiber/v2"
)

func ServeIndex(api *types.Api) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		routeParts := utils.ExtractRouteParts(ctx.Path())
		fmt.Println(routeParts, len(routeParts))
		if len(routeParts) < 2 {
			// Handle non-admin scenario
			return fiber.ErrNotFound
		}

		adminParam := routeParts[1]
		if adminParam != api.AdminSlug {
			// Handle non-admin scenario
			return fiber.ErrNotFound
		}

		return ctx.SendStatus(fiber.StatusNotFound)
	}
}
