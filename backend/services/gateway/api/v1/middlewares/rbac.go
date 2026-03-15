package middlewares

import (
	"github.com/chibx/vuecom/backend/services/gateway/api/v1/response"
	"github.com/chibx/vuecom/backend/services/gateway/internal/constants"
	"github.com/chibx/vuecom/backend/shared/rbac"
	"github.com/gofiber/fiber/v2"
)

func HasPermission(perms ...string) fiber.Handler {
	err401 := fiber.NewError(fiber.StatusUnauthorized, "You do not have the permission to proceed.")
	return func(ctx *fiber.Ctx) error {
		role, _ := ctx.Locals(constants.RoleCtxKey).(*rbac.Role)
		if role == nil {
			return response.FromFiberError(ctx, err401)
		}

		if !role.Has(perms...) {
			return response.FromFiberError(ctx, err401)
		}

		return ctx.Next()
	}
}
