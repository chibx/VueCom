package middlewares

import (
	"errors"
	"net/url"
	"strings"
	"vuecom/gateway/internal/auth"
	"vuecom/gateway/internal/cache"
	"vuecom/gateway/internal/types"
	"vuecom/gateway/internal/types/constants"
	"vuecom/gateway/internal/utils"
	"vuecom/shared/errors/server"
	dbModels "vuecom/shared/models/db"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func getAuthUserFromSession(ctx *fiber.Ctx, api *types.Api, backendUserSess *dbModels.BackendSession) (*dbModels.BackendUser, error) {
	var backendUser *dbModels.BackendUser
	validationErr := auth.ValidateBackendUserSess(ctx, backendUserSess)
	if validationErr != nil {
		var sessionErr *server.SessionErr

		if errors.As(validationErr, &sessionErr) {
			if sessionErr.Type == server.SessionExpired {
				ctx.ClearCookie(constants.BackendCookieKey)
				return nil, server.NewServerErr(fiber.StatusBadRequest, "Session token has expired. Please log in again.")
			}
		}

		return nil, server.NewServerErr(fiber.StatusUnauthorized, "Invalid session")
	}

	if backendUserSess != nil {
		var err error
		backendUser, err = cache.GetBackendUserById(api, int(backendUserSess.UserId), ctx.Context())
		if err != nil {
			return nil, err
		}
	}
	return backendUser, nil
}

func AuthMiddleware(api *types.Api) fiber.Handler {
	logger := api.Deps.Logger
	return func(ctx *fiber.Ctx) error {
		var backendUserSess *dbModels.BackendSession
		var apiKeyData *dbModels.ApiKey
		var backendUser *dbModels.BackendUser
		var tokenErr error
		var authHeader = ctx.Get("Authorization")
		var tokenStr = strings.TrimPrefix(authHeader, "Bearer ")
		var backendToken = ctx.Cookies(constants.BackendCookieKey)

		if tokenStr != "" {
			// TODO: Use tokenStr to validate the api (key) token
			_ = tokenStr
			_ = apiKeyData
		} else if backendToken != "" {
			backendUserSess, tokenErr = cache.GetBackendUserSession(backendToken, api, ctx.Context())
			var authErr error
			backendUser, authErr = getAuthUserFromSession(ctx, api, backendUserSess)
			if authErr != nil {
				logger.Error("failed to get user data from session", zap.Error(authErr))
				// return authErr
			}
		}

		routeParts := utils.ExtractRouteParts(ctx.Path())

		// Validate the user if he is accessing the admin panel
		if len(routeParts) > 1 && routeParts[1] == api.AdminSlug {

			if len(routeParts) > 2 && routeParts[2] == "login" {
				// return ctx.Next() // Skip auth for login page
				return utils.ServeIndex(ctx)
			}

			if strings.TrimSpace(backendToken) == "" {
				logger.Info("Redirecting to login", zap.String("route", routeParts[1]))
				return ctx.Redirect("/" + routeParts[1] + "/login")
			}

			if tokenErr != nil {
				var asTokenErr *server.ServerErr

				if errors.As(tokenErr, &asTokenErr) {
					if asTokenErr.Code == fiber.StatusUnauthorized {
						absoluteUrl := utils.GetAbsoluteUrl(ctx)

						return ctx.Redirect("/"+routeParts[1]+"/login?redirectTo="+url.QueryEscape(absoluteUrl), fiber.StatusSeeOther)
					}

					// Handle other token errors if needed
					return ctx.Status(asTokenErr.Code).SendString(asTokenErr.Message)
				}
			}

			// validationErr := auth.ValidateBackendUserSess(ctx, backendUserSess)
			// if validationErr != nil {
			// 	var sessionErr *serverErrors.SessionErr

			// 	if errors.As(validationErr, &sessionErr) {
			// 		if sessionErr.Type == serverErrors.SessionExpired {
			// 			return ctx.Status(fiber.StatusBadRequest).SendString("Session token has expired. Please log in again.")
			// 		}
			// 	}

			// 	return ctx.Status(fiber.StatusUnauthorized).SendString("Invalid session")
			// }

			// if backendUserSess != nil {
			// 	var err error
			// 	backendUser, err = cache.GetBackendUserById(api, int(backendUserSess.UserId), ctx.Context())
			// 	if err != nil {
			// 		return ctx.Status(fiber.StatusInternalServerError).SendString("Error fetching user data")
			// 	}
			// }

			// Store backend user data in context for downstream handlers
			// return ctx.Next()
		}

		ctx.Locals(constants.ApiKeyCtxKey, tokenStr)
		ctx.Locals(constants.BackendUserCtxKey, backendUser)
		return ctx.Next()
	}
}
