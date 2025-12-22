package middlewares

import (
	"errors"
	"net/url"
	"strings"
	"vuecom/gateway/internal/cache"
	"vuecom/gateway/internal/utils"
	"vuecom/gateway/internal/v1/types"
	userErrors "vuecom/shared/errors/users"

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
		// tokenStr := strings.TrimPrefix(auth, "Bearer ")

		routeParts := utils.ExtractRouteParts(ctx.Path())

		// Validate the user if he is accessing the admin panel
		if len(routeParts) > 1 && routeParts[1] == api.AdminSlug {
			backend_token := ctx.Cookies("backend_session")

			if len(routeParts) > 2 && routeParts[2] == "login" {
				return ctx.Next() // Skip auth for login page
			}

			if strings.TrimSpace(backend_token) == "" {
				return ctx.Redirect(routeParts[1] + "/login") //
			}

			var tokenErr *userErrors.TokenErr
			backendUserData, err := cache.GetBackendUserSession(backend_token, api, ctx.Context())
			if errors.As(err, &tokenErr) {
				if tokenErr.Code == fiber.StatusUnauthorized {
					absoluteUrl := utils.GetAbsoluteUrl(ctx)

					return ctx.Redirect(routeParts[1]+"/login?redirectTo="+url.QueryEscape(absoluteUrl), fiber.StatusSeeOther)
				}
				// Handle other token errors if needed
				return ctx.Status(tokenErr.Code).SendString(tokenErr.Message)
			}

			// Store backend user data in context for downstream handlers
			ctx.Locals("backend_user", backendUserData)
			return ctx.Next()
		}

		return ctx.Next()
	}
}
