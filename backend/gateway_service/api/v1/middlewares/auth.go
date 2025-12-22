package middlewares

import (
	"errors"
	"strings"
	"vuecom/gateway/internal/cache"
	"vuecom/gateway/internal/utils"
	"vuecom/gateway/internal/v1/types"

	"github.com/gofiber/fiber/v2"
)

// Auth middleware: Validates access token.
// func AuthMiddleware(api *types.Api) fiber.Handler {
// 	return func(c *fiber.Ctx) error {
// 		auth := c.Get("Authorization")
// 		if auth == "" || !strings.HasPrefix(auth, "Bearer ") {
// 			return c.Status(fiber.StatusUnauthorized).SendString("Unauthorized")
// 		}
// 		tokenStr := strings.TrimPrefix(auth, "Bearer ")

// 		claims := jwt.MapClaims{}
// 		token, err := jwt.ParseWithClaims(tokenStr, &claims, func(token *jwt.Token) (any, error) {
// 			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
// 				return nil, fmt.Errorf("unexpected signing method")
// 			}
// 			return api.Config.SecretKey, nil
// 		})
// 		if err != nil || !token.Valid {
// 			return c.Status(fiber.StatusUnauthorized).SendString("Invalid access token")
// 		}

// 		// Optional: If high security, check if user is still valid (e.g., not banned) via DB.
//		c.Locals("user_id", claims["sub"])
// 		return c.Next()
// 	}
// }

func AuthMiddleware(api *types.Api) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		auth := ctx.Get("Authorization")
		if auth == "" || !strings.HasPrefix(auth, "Bearer ") {
			return ctx.Status(fiber.StatusUnauthorized).SendString("Unauthorized")
		}
		tokenStr := strings.TrimPrefix(auth, "Bearer ")

		routeParts := utils.ExtractRouteParts(ctx.Path())

		if len(routeParts) > 1 {
			if routeParts[1] == api.AdminSlug {
				backend_token := ctx.Cookies("backend_session")
				if strings.TrimSpace(backend_token) == "" {
					return ctx.Redirect(routeParts[1] + "/login") //
				}

				_, err := cache.GetBackendUserSession(backend_token, api, ctx.Context())
				if errors.As() {
					return ctx.Status().Redirect(routeParts[1] + "/login")
				}
			}
		}

		return ctx.Next()
	}
}
