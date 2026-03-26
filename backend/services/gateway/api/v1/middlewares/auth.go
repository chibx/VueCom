package middlewares

import (
	"errors"
	"strings"

	"github.com/chibx/vuecom/backend/shared/errors/server"
	userModels "github.com/chibx/vuecom/backend/shared/models/db/users"
	reqctx "github.com/chibx/vuecom/backend/shared/reqctx"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"

	"github.com/chibx/vuecom/backend/services/gateway/internal/auth"
	"github.com/chibx/vuecom/backend/services/gateway/internal/cache"
	"github.com/chibx/vuecom/backend/services/gateway/internal/constants"
	"github.com/chibx/vuecom/backend/services/gateway/internal/global"
	"github.com/chibx/vuecom/backend/services/gateway/internal/types"
	"github.com/chibx/vuecom/backend/services/gateway/internal/utils"

	"github.com/gofiber/fiber/v2"
)

func getAuthUserFromSession(ctx *fiber.Ctx, api *types.Api, backendUserSess *userModels.BackendSession) (*userModels.BackendUser, error) {
	var backendUser *userModels.BackendUser
	validationErr := auth.ValidateBackendUserSess(ctx, backendUserSess)
	if validationErr != nil {
		var sessionErr *server.SessionErr

		if errors.As(validationErr, &sessionErr) {
			if sessionErr.Type == server.SessionExpired {
				ctx.ClearCookie(constants.BackendRefreshTkKey)
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
	logger := global.Logger()
	return func(c *fiber.Ctx) error {
		// var backendUserSess *userModels.BackendSession
		var apiKeyData *userModels.ApiKey
		var backendUser *reqctx.BackendUser
		// var tokenErr error
		var authHeader = c.Get("Authorization")
		var tokenStr = strings.TrimPrefix(authHeader, "Bearer ")
		var backendToken = strings.TrimSpace(c.Cookies(constants.BackendAccessTkKey))
		// var tokenGroup = strings.Split(backendToken, ".")

		if tokenStr != "" {
			// TODO: Use tokenStr to validate the api (key) token
			_ = tokenStr
			_ = apiKeyData

			// This should be the api key struct
			// We will also check for customer login from here
			c.Locals(constants.ApiKeyCtxKey, apiKeyData)
		} else if backendToken != "" {
			//

			// if len(tokenGroup) < 2 {
			// 	// I will just skip
			// 	return ctx.Next()
			// }
			// tokenId := tokenGroup[0]
			// backendUserSess, tokenErr = auth.GetBackendUserSession(ctx.Context(), tokenId, api)
			// if tokenErr != nil {
			// 	logger.Error("failed to get user session from cache", zap.Error(tokenErr))
			// } else {
			// 	var authErr error
			// 	backendUser, authErr = getAuthUserFromSession(ctx, api, backendUserSess)
			// 	if authErr != nil {
			// 		logger.Error("failed to get user data from session", zap.Error(authErr))
			// 	}
			// }

			validJWT, err := auth.ValidateBackendAccessToken(api, backendToken, api.Config.SecretKey)
			if err == nil {
				backendUser = &reqctx.BackendUser{ID: validJWT.UserID}
				userPerm, ok := global.UserPermCache.Get(validJWT.UserID)
				if !ok {
					permSet, err := utils.RefetchRoleCache(c.Context(), api, validJWT.UserID)
					if err != nil {
						global.Logger().Error("Failed to refresh role cache", zap.Error(err))
					}
					userPerm = permSet
				}

				c.Locals(constants.RoleCtxKey, userPerm)
				c.Locals(constants.BackendUserCtxKey, backendUser)
			} else {
				logger.Error("Error during authentication", zap.Error(err))

				if errors.Is(err, jwt.ErrTokenExpired) {
					// TODO: Add some kind of warning
				}
			}
		}

		return c.Next()
	}
}
