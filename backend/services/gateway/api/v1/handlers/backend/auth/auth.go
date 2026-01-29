package auth

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/chibx/vuecom/backend/services/gateway/api/v1/response"
	"github.com/chibx/vuecom/backend/services/gateway/internal/auth"
	"github.com/chibx/vuecom/backend/services/gateway/internal/constants"
	"github.com/chibx/vuecom/backend/services/gateway/internal/dto"
	"github.com/chibx/vuecom/backend/services/gateway/internal/types"
	"github.com/chibx/vuecom/backend/services/gateway/internal/utils"
	userModels "github.com/chibx/vuecom/backend/shared/models/db/users"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func Register(api *types.Api) fiber.Handler {
	logger := utils.Logger()
	return func(ctx *fiber.Ctx) error {
		var err error
		errLogin500 := fiber.NewError(fiber.StatusInternalServerError, "Error occurred while logging you in, please try again")
		form, err := ctx.MultipartForm()
		if err != nil {
			logger.Error("Error reading multipart for route " + ctx.Path())
			return response.FromFiberError(ctx, errLogin500)
		}

		value := form.Value
		fmt.Println(value)

		return nil
	}
}

func Login(api *types.Api) fiber.Handler {
	logger := utils.Logger()
	db := api.Deps.DB
	return func(ctx *fiber.Ctx) error {
		var backendUser *dto.UserForLogin
		var err error
		errLogin500 := fiber.NewError(fiber.StatusInternalServerError, "Error occurred while logging you in, please try again")

		password := strings.TrimSpace(ctx.FormValue("password"))
		username := strings.TrimSpace(ctx.FormValue("username"))

		if len(username) > constants.MaxUsernameLimit {
			return response.WriteResponse(ctx, fiber.StatusBadRequest, "Please enter a password less than "+strconv.Itoa(constants.MaxUsernameLimit)+" characters")
		}

		if len(username) == 0 {
			return response.WriteResponse(ctx, fiber.StatusBadRequest, "Please enter a character as your username")
		}

		if len(password) > constants.MaxPasswordLimit {
			return response.WriteResponse(ctx, fiber.StatusBadRequest, "Please enter a password less than "+strconv.Itoa(constants.MaxPasswordLimit)+" characters")
		}

		if len(password) == 0 {
			return response.WriteResponse(ctx, fiber.StatusBadRequest, "Please enter valid character as your password")
		}

		backendUser, err = db.BackendUsers().GetUserByNameForLogin(ctx.Context(), username)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return response.WriteResponse(ctx, fiber.StatusUnauthorized, "Invalid username and/or password")
			}

			logger.Error("Database error during backend login", zap.Error(err))
			return response.FromFiberError(ctx, errLogin500)
		}

		match, err := auth.CompareRawAndHash(password, backendUser.PasswordHash)

		if err != nil {
			logger.Error("Error verifying backend login password", zap.Error(err))
			return response.FromFiberError(ctx, errLogin500)
		}

		if !match {
			return response.WriteResponse(ctx, fiber.StatusUnauthorized, "Invalid username and/or password")
		}

		var refreshTokenExp = time.Now().Add(constants.BackendRefreshTkDur)
		var accessTokenExp = time.Now().Add(constants.BackendAccessTkDur)
		var deviceId = ctx.Cookies(constants.DeviceIDKey)
		var ipAddr = ctx.IP()
		var compositeRefreshToken = ctx.Cookies(constants.BackendRefreshTkKey)
		if deviceId == "" {
			deviceUUID, err := uuid.NewRandom()
			if err != nil {
				logger.Error("Error generating deviceUUID", zap.Error(err))
				return response.FromFiberError(ctx, errLogin500)
			}

			deviceId = deviceUUID.String()
			if deviceId == "" {
				logger.Error("Invalid uuid string")
				return response.FromFiberError(ctx, errLogin500)
			}

			// Set the deviceId anyways
			ctx.Cookie(&fiber.Cookie{
				Name:     constants.DeviceIDKey,
				Value:    deviceId,
				Expires:  time.Now().Add(constants.DeviceIDDur),
				SameSite: "Strict",
				HTTPOnly: true,
				Secure:   true,
			})
		}

		compositeRefreshToken, refreshTokenHash, err := auth.CompositeRefreshToken()
		if err != nil {
			logger.Error("Error generating composite refresh token", zap.Error(err))
			return response.FromFiberError(ctx, errLogin500)
		}

		backendSession := userModels.BackendSession{
			UserId:           backendUser.ID,
			RefreshTokenHash: refreshTokenHash,
			LastIP:           ipAddr,
			DeviceId:         deviceId,
			ExpiresAt:        refreshTokenExp,
		}
		err = auth.CreateBackendSession(ctx.Context(), &backendSession, api)
		if err != nil {
			return response.FromFiberError(ctx, errLogin500)
		}

		accessToken, err := auth.GenerateBackendAccessToken(api, int(backendUser.ID))
		if err != nil {
			return response.FromFiberError(ctx, errLogin500)
		}

		ctx.Cookie(&fiber.Cookie{
			Name:     constants.BackendRefreshTkKey,
			Value:    compositeRefreshToken,
			Expires:  refreshTokenExp,
			SameSite: "Strict",
			HTTPOnly: true,
			Secure:   true,
		})

		ctx.Cookie(&fiber.Cookie{
			Name:     constants.BackendAccessTkKey,
			Value:    accessToken,
			Expires:  accessTokenExp,
			SameSite: "Strict",
			HTTPOnly: true,
			Secure:   true,
		})

		return nil
	}
}
