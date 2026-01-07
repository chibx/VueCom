package middlewares

import (
	"net/url"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"vuecom/gateway/internal/constants"
	"vuecom/gateway/internal/types"
	"vuecom/gateway/internal/utils"
	dbModels "vuecom/shared/models/db"

	"github.com/gofiber/fiber/v2"
)

// func ServeIndex(api *types.Api) fiber.Handler {
// 	return func(ctx *fiber.Ctx) error {
// 		logger := api.Deps.Logger
// 		routeParts := utils.ExtractRouteParts(ctx.Path())
// 		var backendToken = strings.TrimSpace(ctx.Cookies(constants.BackendCookieKey))
// 		var backendUser, _ = ctx.Locals(constants.BackendUserCtxKey).(*dbModels.BackendUser)

// 		// Validate the user if he is accessing the admin panel
// 		if len(routeParts) > 1 && routeParts[1] == api.AdminSlug {

// 			if len(routeParts) > 2 && routeParts[2] == "login" {
// 				if backendUser != nil {
// 					return ctx.Redirect("/" + routeParts[1] + "/dashboard")
// 				}
// 				return utils.ServeIndex(ctx)
// 			}

// 			absoluteUrl := utils.GetAbsoluteUrl(ctx)
// 			if backendToken == "" {
// 				logger.Info("Redirecting to login", zap.String("route", routeParts[1]))
// 				return ctx.Redirect("/" + routeParts[1] + "/login")
// 			}

// 			if backendUser == nil {
// 				return ctx.Redirect("/"+routeParts[1]+"/login?redirectTo="+url.QueryEscape(absoluteUrl), fiber.StatusSeeOther)
// 			}

// 		}

// 		return ctx.Next()
// 	}
// }

func ServeIndex(api *types.Api) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		logger := api.Deps.Logger
		absoluteUrl := utils.GetAbsoluteUrl(ctx)
		routeParts := utils.ExtractRouteParts(ctx.Path())
		var backendToken = strings.TrimSpace(ctx.Cookies(constants.BackendCookieKey))
		var backendUser, _ = ctx.Locals(constants.BackendUserCtxKey).(*dbModels.BackendUser)
		var isLoginRoute = len(routeParts) == 2 && routeParts[1] == "login"
		var redirectTo string = "?redirectTo=" + url.QueryEscape(absoluteUrl)
		// Validate the user if he is accessing the admin panel
		// if len(routeParts) > 1 && routeParts[1] == api.AdminSlug {

		if isLoginRoute {
			if backendUser != nil {
				return ctx.Redirect("/dashboard")
			}
			// return utils.ServeIndex(ctx)
			return ctx.Next()
		}

		// prevent redirect on api route
		if len(routeParts) == 0 || (len(routeParts) > 1 && routeParts[1] != "api") {
			if backendToken == "" {
				logger.Info("Redirecting to login")
				return ctx.Redirect("/login" + redirectTo)
			}

			if backendUser == nil {
				logger.Info("Redirecting to login")
				return ctx.Redirect("/login"+redirectTo, fiber.StatusSeeOther)
			}
		}

		// }

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
