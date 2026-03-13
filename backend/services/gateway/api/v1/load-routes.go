package v1

import (
	"errors"
	"strings"
	"time"

	"github.com/chibx/vuecom/backend/services/gateway/api/v1/handlers/backend"
	"github.com/chibx/vuecom/backend/services/gateway/api/v1/handlers/customer"
	"github.com/chibx/vuecom/backend/services/gateway/api/v1/middlewares"
	"github.com/chibx/vuecom/backend/services/gateway/api/v1/response"
	serverErrors "github.com/chibx/vuecom/backend/shared/errors/server"
	"go.uber.org/zap"

	// "github.com/chibx/vuecom/backend/services/gateway/internal/constants"
	"github.com/chibx/vuecom/backend/services/gateway/internal/constants"
	"github.com/chibx/vuecom/backend/services/gateway/internal/types"
	"github.com/chibx/vuecom/backend/services/gateway/internal/utils"

	"github.com/gofiber/fiber/v2"
)

func HandleRegisterRoute(api *types.Api) fiber.Handler {
	db := api.Deps.DB
	logger := utils.Logger()

	return func(ctx *fiber.Ctx) error {
		token := strings.TrimSpace(ctx.Params("token"))
		if token == "" {
			return response.FromFiberError(ctx, fiber.ErrForbidden)
		}

		signupTk, err := db.BackendUsers().GetRegToken(ctx.Context(), token)

		if err != nil {
			if !errors.Is(err, serverErrors.ErrDBRecordNotFound) {
				logger.Error("Error checking for signup token", zap.Error(err))
			}

			return response.FromFiberError(ctx, fiber.ErrForbidden)
		}

		if time.Now().After(signupTk.ExpiryAt) {
			err := db.BackendUsers().DeleteRegToken(ctx.Context(), signupTk.Token)
			if err != nil {
				logger.Error("Error deleting expired reg token")
			}
			return response.WriteResponse(ctx, fiber.StatusBadRequest, "Expired Registration Link.")
		}

		return ctx.Next()
	}
}

// Potentially Long Function | Just stack all the routes in here
func LoadRoutes(app fiber.Router, api *types.Api) {
	// app.Use(middlewares.AuthMiddleware(api), middlewares.ServeIndex(api), middlewares.ServeAssets())
	app.Get("/api/health", func(ctx *fiber.Ctx) error {
		return ctx.Status(fiber.StatusOK).SendString("OK")
	})
	app.Use(middlewares.ServeAssets(), middlewares.AuthMiddleware(api), middlewares.RedirectCommon(api))

	backend.LoadRoutes(app, api)
	customer.LoadRoutes(app, api)

	app.Use(middlewares.ServeIndex(api))

	// app.Static("*", "./"+constants.PublicFolder, fiber.Static{})
	app.Get("/register/:token", HandleRegisterRoute(api))
	app.Get("*", func(ctx *fiber.Ctx) error {
		return ctx.SendFile("./" + constants.PublicFolder + "/index.html")
	})
}
