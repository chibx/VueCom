package auth

import (
	"fmt"

	"github.com/chibx/vuecom/backend/services/gateway/internal/types"
	"github.com/chibx/vuecom/backend/services/gateway/internal/utils"
	"github.com/gofiber/fiber/v2"
)

func Login(api *types.Api) fiber.Handler {
	logger := utils.Logger()
	return func(ctx *fiber.Ctx) error {
		form, err := ctx.MultipartForm()
		if err != nil {
			logger.Error("Error reading multipart for route " + ctx.Path())
		}

		value := form.Value
		fmt.Println(value)

		return nil
	}
}

func Register(api *types.Api) fiber.Handler {
	logger := utils.Logger()
	return func(ctx *fiber.Ctx) error {
		form, err := ctx.MultipartForm()
		if err != nil {
			logger.Error("Error reading multipart for route " + ctx.Path())
		}

		value := form.Value
		fmt.Println(value)

		return nil
	}
}
