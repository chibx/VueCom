package admin

import (
	"errors"

	"github.com/chibx/vuecom/backend/services/gateway/api/v1/middlewares"
	"github.com/chibx/vuecom/backend/services/gateway/api/v1/response"
	"github.com/chibx/vuecom/backend/services/gateway/internal/types"
	"github.com/chibx/vuecom/backend/services/gateway/internal/utils"

	serverErrors "github.com/chibx/vuecom/backend/shared/errors/server"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func RegisterRoutes(app fiber.Router, api *types.Api) {
	appGroup := app.Group("/api/app", middlewares.BackendRateLimit(api))
	appGroup.Post("/initialize", InitializeApp(api))
	appGroup.Post("/create-owner", RegisterOwner(api))

	appGroup.Get("/admin-exist", func(ctx *fiber.Ctx) error {
		logger := utils.Logger()
		exists, err := DoesOwnerExist(ctx, api)

		if err != nil {
			logger.Error("Error checking for existing users", zap.Error(err))
			if errors.Is(err, serverErrors.ErrDBRecordNotFound) {
				return response.WriteResponse(ctx, fiber.StatusBadRequest, "Owner does not exist")
			}
			return response.WriteResponse(ctx, fiber.StatusInternalServerError, "An Error occurred, please try again")
		}

		return response.WriteResponse(ctx, fiber.StatusOK, "Success", fiber.Map{
			"exists": exists,
		})
	})
}
