package middlewares

import (
	"github.com/chibx/vuecom/backend/services/gateway/api/v1/response"
	"github.com/chibx/vuecom/backend/services/gateway/internal/constants"
	"github.com/chibx/vuecom/backend/services/gateway/internal/types"
	"github.com/chibx/vuecom/backend/services/gateway/internal/utils"
	userModels "github.com/chibx/vuecom/backend/shared/models/db/users"
	"github.com/gofiber/fiber/v2"
)

// I dont know what to call this function honestly
func RedirectCommon(api *types.Api) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		routeParts := utils.ExtractRouteParts(ctx.Path())
		var backendUser, _ = ctx.Locals(constants.BackendUserCtxKey).(*userModels.BackendUser)
		var isApiRoute = len(routeParts) > 1 && routeParts[1] == "api"
		var isAppInitRoute = len(routeParts) == 3 && routeParts[1] == "app" && routeParts[2] == "init"
		var isAdminCreateRoute = len(routeParts) == 3 && routeParts[1] == "app" && routeParts[2] == "create-user"

		// TODO: I will refactor these checks later
		// ---------------------------------------------
		// I might choose to handle this client-side instead as the api route meant to do the work will be the one guarded
		if isAppInitRoute && api.IsAppInit {
			if backendUser != nil {
				return ctx.Redirect("/dashboard")
			} else {
				return ctx.Redirect("/login")
			}
		}

		// I might choose to handle this client-side instead as the api route meant to do the work will be the one guarded
		if isAdminCreateRoute && api.HasAdmin {
			if backendUser != nil {
				return ctx.Redirect("/dashboard")
			} else {
				return ctx.Redirect("/login")
			}
		}

		if isAppInitRoute {
			return ctx.Next()
		}

		if !api.IsAppInit {
			if isApiRoute {
				return response.WriteResponse(ctx, fiber.StatusBadRequest, "There is no app initialized!!!")
			}
			return ctx.Redirect("/app/init")
		}

		if isAdminCreateRoute {
			return ctx.Next()
		}

		if !api.HasAdmin {
			if isApiRoute {
				return response.WriteResponse(ctx, fiber.StatusBadRequest, "There is no app owner!!!")
			}
			return ctx.Redirect("/app/create-user")
		}

		return ctx.Next()
	}
}
