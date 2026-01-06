package middlewares

import (
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"vuecom/gateway/internal/types"
	"vuecom/gateway/internal/types/constants"
	"vuecom/gateway/internal/utils"
	dbModels "vuecom/shared/models/db"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func ServeIndex(api *types.Api) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		logger := api.Deps.Logger
		routeParts := utils.ExtractRouteParts(ctx.Path())
		var backendToken = strings.TrimSpace(ctx.Cookies(constants.BackendCookieKey))
		var backendUser, _ = ctx.Locals(constants.BackendUserCtxKey).(*dbModels.BackendUser)

		// Validate the user if he is accessing the admin panel
		if len(routeParts) > 1 && routeParts[1] == api.AdminSlug {

			if len(routeParts) > 2 && routeParts[2] == "login" {
				if backendUser != nil {
					return ctx.Redirect("/" + routeParts[1] + "/dashboard")
				}
				return utils.ServeIndex(ctx)
			}

			if backendToken == "" {
				logger.Info("Redirecting to login", zap.String("route", routeParts[1]))
				return ctx.Redirect("/" + routeParts[1] + "/login")
			}

			if backendUser == nil {
				absoluteUrl := utils.GetAbsoluteUrl(ctx)
				return ctx.Redirect("/"+routeParts[1]+"/login?redirectTo="+url.QueryEscape(absoluteUrl), fiber.StatusSeeOther)
			}

		}

		return ctx.Next()
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
