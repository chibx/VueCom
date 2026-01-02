package middlewares

import (
	"vuecom/gateway/api/v1/request"
	"vuecom/gateway/internal/types"

	"github.com/gofiber/fiber/v2"
)

// Redirects the user to login in the appropriate scenario (if needed)
func RequireLogin(api *types.Api) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		// backend_token := ctx.Cookies(request.BACKEND_TOKEN)
		_ = ctx.Cookies(request.CUSTOMER_TOKEN)
		_ = ctx.Cookies(request.BACKEND_TOKEN)

		return ctx.Next()
	}
}
