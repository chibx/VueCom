package middlewares

import (
	"errors"
	"strings"

	"github.com/chibx/vuecom/backend/shared/errors/server"
	userModels "github.com/chibx/vuecom/backend/shared/models/db/users"

	"github.com/chibx/vuecom/backend/services/gateway/internal/auth"
	"github.com/chibx/vuecom/backend/services/gateway/internal/cache"
	"github.com/chibx/vuecom/backend/services/gateway/internal/constants"
	"github.com/chibx/vuecom/backend/services/gateway/internal/types"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func getAuthUserFromSession(ctx *fiber.Ctx, api *types.Api, backendUserSess *userModels.BackendSession) (*userModels.BackendUser, error) {
	var backendUser *userModels.BackendUser
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
		var backendUserSess *userModels.BackendSession
		var apiKeyData *userModels.ApiKey
		var backendUser *userModels.BackendUser
		var tokenErr error
		var authHeader = ctx.Get("Authorization")
		var tokenStr = strings.TrimPrefix(authHeader, "Bearer ")
		var backendToken = strings.TrimSpace(ctx.Cookies(constants.BackendCookieKey))
		var tokenGroup = strings.Split(backendToken, ".")

		if tokenStr != "" {
			// TODO: Use tokenStr to validate the api (key) token
			_ = tokenStr
			_ = apiKeyData

			// This should be the api key struct
			ctx.Locals(constants.ApiKeyCtxKey, apiKeyData)
		}
		// else
		if backendToken != "" {
			//

			if len(tokenGroup) < 2 {
				// I will just skip
				return ctx.Next()
			}

			tokenId := tokenGroup[0]
			backendUserSess, tokenErr = auth.GetBackendUserSession(tokenId, api, ctx.Context())
			if tokenErr != nil {
				logger.Error("failed to get user session from cache", zap.Error(tokenErr))
			} else {
				var authErr error
				backendUser, authErr = getAuthUserFromSession(ctx, api, backendUserSess)
				if authErr != nil {
					logger.Error("failed to get user data from session", zap.Error(authErr))
				}
			}

		}

		ctx.Locals(constants.BackendUserCtxKey, backendUser)

		return ctx.Next()
	}
}
