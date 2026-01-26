package auth

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/chibx/vuecom/backend/services/gateway/api/v1/response"
	"github.com/chibx/vuecom/backend/services/gateway/internal/constants"
	"github.com/chibx/vuecom/backend/services/gateway/internal/types"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func Login(api *types.Api) fiber.Handler {
	logger := api.Deps.Logger
	return func(ctx *fiber.Ctx) error {
		form, err := ctx.MultipartForm()
		if err != nil {
			logger.Error("Error reading multipart for route " + ctx.Path())
		}

		value := form.Value
		passwordValues := value["password"]
		usernameValues := value["username"]
		if passwordValues == nil || usernameValues == nil {
			return response.WriteResponse(ctx, fiber.StatusBadRequest, "One or more fields are missing")
		}

		password := passwordValues[0]
		username := usernameValues[0]

		if len(password) > constants.MaxPasswordLimit {
			return response.WriteResponse(ctx, fiber.StatusBadRequest, "Please enter a password less than "+strconv.Itoa(constants.MaxPasswordLimit)+" characters")
		}

		if strings.Contains(password, " ") || len(password) == 0 {
			return response.WriteResponse(ctx, fiber.StatusBadRequest, "Please enter valid character. Spaces aren't allowed")
		}

		if len(username) > constants.MaxUsernameLimit {
			return response.WriteResponse(ctx, fiber.StatusBadRequest, "Please enter a password less than "+strconv.Itoa(constants.MaxUsernameLimit)+" characters")
		}

		if len(password) == 0 {
			return response.WriteResponse(ctx, fiber.StatusBadRequest, "Please enter a character")
		}

		logger.Info("Form fields", zap.Strings("fields", []string{username, password}))

		// auth.ComparePasswordAndHash()

		return nil
	}
}

func Register(api *types.Api) fiber.Handler {
	logger := api.Deps.Logger
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
