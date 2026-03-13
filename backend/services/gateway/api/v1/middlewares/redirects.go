package middlewares

import (
	"github.com/chibx/vuecom/backend/services/gateway/api/v1/response"
	"github.com/chibx/vuecom/backend/services/gateway/internal/constants"
	"github.com/chibx/vuecom/backend/services/gateway/internal/types"
	"github.com/chibx/vuecom/backend/services/gateway/internal/utils"
	reqctx "github.com/chibx/vuecom/backend/shared/reqctx"
	"github.com/gofiber/fiber/v2"
)

// I dont know what to call this function honestly
func RedirectCommon(api *types.Api) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		path := utils.WithTrailingSlash(ctx.Path())
		routeParts := utils.ExtractRouteParts(path)
		var backendUser, _ = ctx.Locals(constants.BackendUserCtxKey).(*reqctx.BackendUser)
		var isApiRoute = len(routeParts) > 1 && routeParts[1] == "api"
		var isAppInitPage = path == "/app/initialize/"
		var isAdminCreatePage = path == "/app/create-owner/"
		var isAppInitRoute = path == "/api/app/initialize/"
		var isAdminCreateRoute = path == "/api/app/create-owner/"

		// TODO: I will refactor these checks later
		if isAppInitRoute || isAdminCreateRoute {
			return ctx.Next()
		}

		// ---------------------------------------------
		// I might choose to handle this client-side instead as the api route meant to do the work will be the one guarded
		if isAppInitPage && api.IsAppInit {
			if backendUser != nil {
				return ctx.Redirect("/dashboard")
			} else {
				return ctx.Redirect("/login")
			}
		}

		// I might choose to handle this client-side instead as the api route meant to do the work will be the one guarded
		if isAdminCreatePage && api.HasAdmin {
			if backendUser != nil {
				return ctx.Redirect("/dashboard")
			} else {
				return ctx.Redirect("/login")
			}
		}

		if isAppInitPage {
			return ctx.Next()
		}

		if !api.IsAppInit {
			if isApiRoute {
				return response.WriteResponse(ctx, fiber.StatusBadRequest, "There is no app initialized!!!")
			}
			return ctx.Redirect("/app/initialize")
		}

		if isAdminCreatePage {
			return ctx.Next()
		}

		if !api.HasAdmin {
			if isApiRoute {
				return response.WriteResponse(ctx, fiber.StatusBadRequest, "There is no app owner!!!")
			}
			return ctx.Redirect("/app/create-owner")
		}

		return ctx.Next()
	}
}
