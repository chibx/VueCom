package auth

import (
	"crypto/subtle"
	"errors"
	"strconv"
	"strings"
	"time"

	backendusers "github.com/chibx/vuecom/backend/services/gateway/api/v1/request/backend_users"
	"github.com/chibx/vuecom/backend/services/gateway/api/v1/response"
	"github.com/chibx/vuecom/backend/services/gateway/internal/auth"
	"github.com/chibx/vuecom/backend/services/gateway/internal/constants"
	"github.com/chibx/vuecom/backend/services/gateway/internal/dto"
	"github.com/chibx/vuecom/backend/services/gateway/internal/types"
	"github.com/chibx/vuecom/backend/services/gateway/internal/utils"
	userModels "github.com/chibx/vuecom/backend/shared/models/db/users"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"

	serverErrors "github.com/chibx/vuecom/backend/shared/errors/server"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

// Handles the registration links in the url
func Register(api *types.Api) fiber.Handler {
	logger := utils.Logger()
	db := api.Deps.DB
	errRegister500 := fiber.NewError(fiber.StatusInternalServerError, "Error occurred while creating your account, please try again")
	return func(ctx *fiber.Ctx) error {
		var err error
		// var errorBag = []serverErrors.ErrorDetail{}
		var userForRegister = new(backendusers.CreateBackendUserRequest)
		var regTokenJWT = ctx.Cookies("reg_token")
		var now = time.Now()
		if len(strings.TrimSpace(regTokenJWT)) == 0 {
			return response.WriteResponse(ctx, fiber.StatusBadRequest, "Invalid Request!")
		}

		regToken, err := auth.ValidateRegToken(api, regTokenJWT, api.Config.SecretKey)
		if err != nil {
			var serverErrors = new(serverErrors.ServerErr)
			if errors.As(err, &serverErrors) {
				return response.WriteResponse(ctx, fiber.StatusBadRequest, serverErrors.Message)
			}
			return response.FromFiberError(ctx, errRegister500)
		}

		if now.After(regToken.ExpiresAt.Time) {
			logger.Warn("Reg Token: Used after jwt expiration")
			return response.WriteResponse(ctx, fiber.StatusBadRequest, "Expired Registration Token!!")
		}

		err = ctx.BodyParser(userForRegister)
		if err != nil {
			logger.Error("Error occured while parsing login values", zap.Error(err))
			return response.FromFiberError(ctx, errRegister500)
		}

		err = utils.Validator().Struct(userForRegister)

		if err != nil {
			if errors.Is(err, &validator.InvalidValidationError{}) {
				logger.Error("InvalidValidationError while registering a user", zap.Error(err))
				return response.WriteResponse(ctx, fiber.ErrBadRequest.Code, errRegister500.Message)
			}

			errorBag := serverErrors.ValErrToBag(err)

			if len(errorBag) > 0 {
				return response.WriteResponse(ctx, fiber.StatusBadRequest, "One or more fields are invalid", errorBag)
			}
		}

		tokenStruc, err := db.BackendUsers().GetRegToken(ctx.Context(), regTokenJWT)

		if err != nil {
			if errors.Is(err, serverErrors.ErrDBRecordNotFound) {
				return response.WriteResponse(ctx, fiber.StatusBadRequest, "Registration Token not found.")
			}

			return response.FromFiberError(ctx, errRegister500)
		}

		if now.After(tokenStruc.ExpiryAt) {
			logger.Warn("Reg Token: Used after db expiration")
			return response.WriteResponse(ctx, fiber.StatusBadRequest, "Expired Registration Token!!")
		}

		if subtle.ConstantTimeCompare([]byte(tokenStruc.Code), []byte(userForRegister.Code)) == 0 {
			logger.Warn("Invalid Code used")
			// TODO: Maybe add a counter that would delete the token and alert the app owner of a potential cyber attack
			return response.WriteResponse(ctx, fiber.StatusUnauthorized, "You cannot proceed")
		}

		if len(userForRegister.UserName) > constants.MaxUsernameLimit {
			return response.WriteResponse(ctx, fiber.StatusBadRequest, "Please enter a username less than "+strconv.Itoa(constants.MaxUsernameLimit)+" characters")
		}

		if len(userForRegister.Password) > constants.MaxPasswordLimit {
			return response.WriteResponse(ctx, fiber.StatusBadRequest, "Please enter a password less than "+strconv.Itoa(constants.MaxPasswordLimit)+" characters")
		}

		user, err := userForRegister.ToDBBackendUser(ctx.Context(), api, ctx)
		if err != nil {
			var serverErr = new(serverErrors.ServerErr)
			if errors.As(err, &serverErr) {
				logger.Error(serverErr.Message)
				return response.WriteResponse(ctx, serverErr.Code, serverErr.Message)
			}

			logger.Error("Error converting user to db backend user", zap.Error(err))
			return response.FromFiberError(ctx, errRegister500)
		}

		err = db.BackendUsers().CreateUser(ctx.Context(), user)
		if err != nil {
			logger.Error("DB Error creating backend user", zap.Error(err))
			return response.FromFiberError(ctx, errRegister500)
		}

		return response.WriteResponse(ctx, fiber.StatusOK, "User created successfully.")
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
			return response.WriteResponse(ctx, fiber.StatusBadRequest, "Please enter a username less than "+strconv.Itoa(constants.MaxUsernameLimit)+" characters")
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
			if errors.Is(err, serverErrors.ErrDBRecordNotFound) {
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
		// var refreshToken = ctx.Cookies(constants.BackendRefreshTkKey)
		// db.BackendUsers().DeleteSession(ctx.Context(), &userModels.BackendSession{
		// 	DeviceId: deviceId,
		// 	RefreshTokenHash: refreshToken,
		// })

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

		refreshToken, refreshTokenHash, err := auth.CompositeRefreshToken()
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
			Value:    refreshToken,
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

func CreateSignupToken(api *types.Api) fiber.Handler {
	_ = utils.Logger()
	_ = api.Deps.DB
	return func(ctx *fiber.Ctx) error {
		return nil
	}
}
