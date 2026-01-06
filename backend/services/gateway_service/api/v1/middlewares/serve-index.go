package middlewares

import (
	"fmt"
	"os"
	"path/filepath"
	"vuecom/gateway/internal/types"
	"vuecom/gateway/internal/types/constants"
	"vuecom/gateway/internal/utils"

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

func ServeAssets() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		// Check if the path exists in the public folder
		path := filepath.Join(constants.PublicFolder, ctx.Path())

		_, err := os.ReadFile(path)

		if err == nil {
			return ctx.SendFile(path)
		}

		return ctx.Next()
	}
}
