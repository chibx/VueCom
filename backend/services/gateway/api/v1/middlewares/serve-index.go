package middlewares

import (
	"net/url"
	"os"
	"path/filepath"
	"slices"
	"strings"

	userModels "github.com/chibx/vuecom/backend/shared/models/db/users"

	"github.com/chibx/vuecom/backend/services/gateway/internal/constants"
	"github.com/chibx/vuecom/backend/services/gateway/internal/types"
	"github.com/chibx/vuecom/backend/services/gateway/internal/utils"

	"github.com/gofiber/fiber/v2"
)

func ServeIndex(api *types.Api) fiber.Handler {
	logger := api.Deps.Logger
	return func(ctx *fiber.Ctx) error {
		absoluteUrl := utils.GetAbsoluteUrl(ctx)
		routeParts := utils.ExtractRouteParts(ctx.Path())
		var backendToken = strings.TrimSpace(ctx.Cookies(constants.BackendRefreshTkKey))
		var backendUser, _ = ctx.Locals(constants.BackendUserCtxKey).(*userModels.BackendUser)
		var isLoginRoute = len(routeParts) == 2 && routeParts[1] == "login"
		var redirectTo = "?redirectTo=" + url.QueryEscape(absoluteUrl)

		if isLoginRoute {
			if backendUser != nil {
				return ctx.Redirect("/dashboard")
			}
			return ctx.Next()
		}

		// prevent redirect on api route
		if len(routeParts) == 1 || (len(routeParts) > 1 && routeParts[1] != "api") {
			if backendToken == "" {
				logger.Info("Redirecting to login")
				return ctx.Redirect("/login"+redirectTo, fiber.StatusSeeOther)
			}

			if backendUser == nil {
				logger.Info("Redirecting to login")
				return ctx.Redirect("/login"+redirectTo, fiber.StatusSeeOther)
			}
		}

		return ctx.Next()
	}
}

func ServeAssets() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		path := ctx.Path()
		routeParts := utils.ExtractRouteParts(path)
		if len(routeParts) > 1 && !slices.Contains(constants.PublicAssets, routeParts[1]) {
			return ctx.Next()
		}
		// Check if the path exists in the public folder
		publicPath := filepath.Join(constants.PublicFolder, path)

		_, err := os.ReadFile(publicPath)

		if err == nil {
			return ctx.SendFile(publicPath)
		}

		return ctx.Next()
	}
}
