package middlewares

import (
	"github.com/chibx/vuecom/backend/services/gateway/api/v1/response"
	"github.com/chibx/vuecom/backend/services/gateway/internal/constants"
	"github.com/chibx/vuecom/backend/shared/rbac"
	"github.com/gofiber/fiber/v2"
)

func HasPermission(perms ...rbac.Permission) fiber.Handler {
	err401 := fiber.NewError(fiber.StatusUnauthorized, "You do not have the permission to proceed.")
	return func(ctx *fiber.Ctx) error {
		permSet, _ := ctx.Locals(constants.RoleCtxKey).(rbac.PermissionSet)
		if permSet == nil {
			return response.FromFiberError(ctx, err401)
		}

		if !permSet.Has(perms...) {
			return response.FromFiberError(ctx, err401)
		}

		return ctx.Next()
	}
}
