package middlewares

import (
	"errors"
	"net/url"
	"strings"
	"vuecom/gateway/internal/auth"
	"vuecom/gateway/internal/cache"
	"vuecom/gateway/internal/utils"
	"vuecom/gateway/internal/v1/types"
	serverErrors "vuecom/shared/errors/server"
	dbModels "vuecom/shared/models/db"

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
		var backendUserSess *dbModels.BackendSession
		var apiKeyData *dbModels.ApiKey
		var backendUser *dbModels.BackendUser
		var tokenErr error

		authHeader := ctx.Get("Authorization")
		if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
			tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
			// TODO: Use tokenStr to validate the api (key) token
			_ = tokenStr
			_ = apiKeyData
		} else {
			backendToken := ctx.Cookies("backend_session")

			backendUserSess, tokenErr = cache.GetBackendUserSession(backendToken, api, ctx.Context())
		}

		routeParts := utils.ExtractRouteParts(ctx.Path())

		// Validate the user if he is accessing the admin panel
		if len(routeParts) > 1 && routeParts[1] == api.AdminSlug {
			backend_token := ctx.Cookies("backend_session")

			if len(routeParts) > 2 && routeParts[2] == "login" {
				// return ctx.Next() // Skip auth for login page
				return utils.ServeIndex(ctx)
			}

			if strings.TrimSpace(backend_token) == "" {
				return ctx.Redirect(routeParts[1] + "/login")
			}

			// var tokenErr *serverErrors.TokenErr
			// backendUserData, err := cache.GetBackendUserSession(backend_token, api, ctx.Context())
			if tokenErr != nil {
				var asTokenErr *serverErrors.TokenErr

				if errors.As(tokenErr, &asTokenErr) {
					if asTokenErr.Code == fiber.StatusUnauthorized {
						absoluteUrl := utils.GetAbsoluteUrl(ctx)

						return ctx.Redirect(routeParts[1]+"/login?redirectTo="+url.QueryEscape(absoluteUrl), fiber.StatusSeeOther)
					}

					// Handle other token errors if needed
					return ctx.Status(asTokenErr.Code).SendString(asTokenErr.Message)
				}
			}

			validationErr := auth.ValidateBackendUserSess(ctx, backendUserSess)
			if validationErr != nil {
				var sessionErr *serverErrors.SessionErr

				if errors.As(validationErr, &sessionErr) {
					if sessionErr.Type == serverErrors.SessionExpired {
						return ctx.Status(fiber.StatusBadRequest).SendString("Session token has expired. Please log in again.")
					}
				}

				return ctx.Status(fiber.StatusUnauthorized).SendString("Invalid session")
			}

			if backendUserSess != nil {
				var err error
				backendUser, err = cache.GetBackendUserById(api, int(backendUserSess.UserId), ctx.Context())
				if err != nil {
					return ctx.Status(fiber.StatusInternalServerError).SendString("Error fetching user data")
				}
			}

			// Store backend user data in context for downstream handlers
			ctx.Locals("backend_user", backendUser)
			return ctx.Next()
		}

		return ctx.Next()
	}
}
