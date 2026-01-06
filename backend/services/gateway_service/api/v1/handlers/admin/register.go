package admin

import (
	"errors"
	"vuecom/gateway/api/v1/response"
	"vuecom/gateway/internal/types"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func RegisterRoutes(app fiber.Router, api *types.Api) {
	app.Post("/initialize-app", func(ctx *fiber.Ctx) error {
		return InitializeApp(ctx, api)
	})
	app.Post("/register-owner", func(ctx *fiber.Ctx) error {
		return RegisterOwner(ctx, api)
	})
	app.Get("/admin-exist", func(ctx *fiber.Ctx) error {
		logger := api.Deps.Logger
		exists, err := DoesOwnerExist(ctx, api)

		if err != nil {
			logger.Error("Error checking for existing users", zap.Error(err))
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return response.NewResponse(ctx, fiber.StatusBadRequest, "Owner does not exist")
			}
			return response.NewResponse(ctx, fiber.StatusInternalServerError, "An Error occurred, please try again")
		}
		return response.NewResponse(ctx, fiber.StatusOK, "", fiber.Map{
			"exists": exists,
		})
	})
}
