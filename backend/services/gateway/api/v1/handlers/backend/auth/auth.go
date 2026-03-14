package auth

import (
	"crypto/sha256"
	"crypto/subtle"
	"errors"
	"strconv"
	"strings"
	"time"

	backendusers "github.com/chibx/vuecom/backend/services/gateway/api/v1/request/backend_users"
	"github.com/chibx/vuecom/backend/services/gateway/api/v1/response"
	backendResp "github.com/chibx/vuecom/backend/services/gateway/api/v1/response/backend_users"
	"github.com/chibx/vuecom/backend/services/gateway/internal/auth"
	"github.com/chibx/vuecom/backend/services/gateway/internal/constants"
	"github.com/chibx/vuecom/backend/services/gateway/internal/dto"
	"github.com/chibx/vuecom/backend/services/gateway/internal/types"
	"github.com/chibx/vuecom/backend/services/gateway/internal/utils"
	"github.com/chibx/vuecom/backend/shared/events"
	userModels "github.com/chibx/vuecom/backend/shared/models/db/users"
	"github.com/google/uuid"
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/hotp"

	serverErrors "github.com/chibx/vuecom/backend/shared/errors/server"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

// Handles the registration links in the url
func Register(api *types.Api) fiber.Handler {
	logger := utils.Logger()
	db := api.Deps.DB
	err500 := fiber.NewError(fiber.StatusInternalServerError, "Error occurred while creating your account, please try again")
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
			return response.FromFiberError(ctx, err500)
		}

		if now.After(regToken.ExpiresAt.Time) {
			logger.Warn("Reg Token: Used after jwt expiration")
			return response.WriteResponse(ctx, fiber.StatusBadRequest, "Expired Registration Token!!")
		}

		err = ctx.BodyParser(userForRegister)
		if err != nil {
			logger.Error("Error occured while parsing login values", zap.Error(err))
			return response.FromFiberError(ctx, err500)
		}

		err = utils.Validator().Struct(userForRegister)

		isFatal, errorBag := serverErrors.HandleValidationError(err)
		if isFatal {
			logger.Error("InvalidValidationError while registering a user", zap.Error(err))
			return response.WriteResponse(ctx, fiber.ErrBadRequest.Code, err500.Message)
		}
		if len(errorBag) > 0 {
			return response.WriteResponse(ctx, fiber.StatusBadRequest, "One or more fields are invalid", errorBag)
		}

		tokenStruc, err := db.BackendUsers().GetRegToken(ctx.Context(), regTokenJWT)

		if err != nil {
			if errors.Is(err, serverErrors.ErrDBRecordNotFound) {
				return response.WriteResponse(ctx, fiber.StatusBadRequest, "Registration Token not found.")
			}

			return response.FromFiberError(ctx, err500)
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
			return response.FromFiberError(ctx, err500)
		}

		err = db.BackendUsers().CreateUser(ctx.Context(), user)
		if err != nil {
			logger.Error("DB Error creating backend user", zap.Error(err))
			return response.FromFiberError(ctx, err500)
		}

		// TODO: Add an ip addition feature

		return response.WriteResponse(ctx, fiber.StatusOK, "User created successfully.")
	}
}

func Login(api *types.Api) fiber.Handler {
	logger := utils.Logger()
	db := api.Deps.DB
	errLogin500 := fiber.NewError(fiber.StatusInternalServerError, "Error occurred while logging you in, please try again")
	return func(ctx *fiber.Ctx) error {
		var backendUser *dto.UserForLogin
		var err error

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

		return response.WriteResponse(ctx, fiber.StatusOK, "Success")
	}
}

func Refresh(api *types.Api) fiber.Handler {
	db := api.Deps.DB
	logger := utils.Logger()
	errLogin500 := fiber.NewError(fiber.StatusInternalServerError, "Couldn't refresh user's session")
	return func(ctx *fiber.Ctx) error {
		refreshTk := ctx.Cookies(constants.BackendRefreshTkKey)
		if refreshTk == "" {
			return response.WriteResponse(ctx, fiber.StatusUnauthorized, "User Session not found. Try logging in again")
		}

		refreshTkHash, err := auth.GenerateHashFromString(refreshTk, auth.DefaultHashParams)
		if err != nil {
			return response.FromFiberError(ctx, errLogin500)
		}
		session, err := db.BackendUsers().GetSessionByTokenHash(ctx.Context(), refreshTkHash)
		if err != nil {
			if errors.Is(err, serverErrors.ErrDBRecordNotFound) {
				return response.WriteResponse(ctx, fiber.StatusUnauthorized, "User Session not found. Try logging in again")
			}

			return response.WriteResponse(ctx, fiber.StatusInternalServerError, "Something went wrong, couldn't refresh user's session.")
		}

		if time.Now().After(session.ExpiresAt) {
			err = db.BackendUsers().DeleteSession(ctx.Context(), session)
			if err != nil {
				logger.Error("Error deleting session token", zap.Error(err))
			}
			return response.WriteResponse(ctx, fiber.StatusUnauthorized, "User Session has expired, you have to login again.")
		}

		var refreshTokenExp = time.Now().Add(constants.BackendRefreshTkDur)
		var accessTokenExp = time.Now().Add(constants.BackendAccessTkDur)
		var deviceId = ctx.Cookies(constants.DeviceIDKey)
		var ipAddr = ctx.IP()

		// TODO: Validate more
		if deviceId != session.DeviceId {
			// This was meant to replace fingerprinting (in a way)
		}

		if ipAddr != session.LastIP {
			// Do some IP range magic or ignore
		}

		if deviceId == "" {
			deviceUUID, err := uuid.NewRandom()
			if err != nil {
				logger.Error("Error generating deviceUUID", zap.Error(err))
				return response.FromFiberError(ctx, errLogin500)
			}

			deviceId = deviceUUID.String()
			if deviceId == "" {
				logger.Error("Invalid uuid string while se")
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
			UserId:           session.UserId,
			RefreshTokenHash: refreshTokenHash,
			LastIP:           ipAddr,
			DeviceId:         deviceId,
			ExpiresAt:        refreshTokenExp,
		}
		err = auth.CreateBackendSession(ctx.Context(), &backendSession, api)
		if err != nil {
			return response.FromFiberError(ctx, errLogin500)
		}

		accessToken, err := auth.GenerateBackendAccessToken(api, int(session.UserId))
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

		return response.WriteResponse(ctx, fiber.StatusOK, "", backendResp.RefreshResp{
			AccessToken: accessToken,
		})
	}
}

func CreateSignupToken(api *types.Api) fiber.Handler {
	logger := utils.Logger()
	db := api.Deps.DB
	err500 := fiber.NewError(fiber.StatusInternalServerError, "Error occurred while creating signup token, please try again")
	return func(ctx *fiber.Ctx) error {
		var reqBody = backendusers.CreateTokenRequest{}
		err := ctx.BodyParser(&reqBody)
		if err != nil {
			return response.FromFiberError(ctx, err500)
		}
		err = utils.Validator().Struct(reqBody)
		isFatal, errorBag := serverErrors.HandleValidationError(err)
		if isFatal {
			logger.Error("InvalidValidationError while creating a signup token", zap.Error(err))
			return response.WriteResponse(ctx, fiber.ErrBadRequest.Code, err500.Message)
		}
		if len(errorBag) > 0 {
			return response.WriteResponse(ctx, fiber.StatusBadRequest, "One or more fields are invalid", errorBag)
		}
		now := time.Now()
		hash := sha256.New()
		hash.Write([]byte(reqBody.Email + strconv.Itoa(int(now.Unix()))))
		final := string(hash.Sum(nil))

		key, err := hotp.Generate(hotp.GenerateOpts{
			AccountName: reqBody.Email,
			Issuer:      api.AppName,
			Digits:      8,
			Algorithm:   otp.AlgorithmSHA256,
			SecretSize:  16,
		})

		if err != nil {
			logger.Error("Error creating code for the registration token.", zap.Error(err))
			return err500
		}

		code := key.String()

		err = db.BackendUsers().CreateRegToken(ctx.Context(), final, reqBody.Supervisor, code)
		if err != nil {
			return response.FromFiberError(ctx, err500)
		}

		// TODO: Implement code sending the token the email
		_, _ = events.DefaultEmitter.Publish(ctx.Context(), &events.Event{
			Type: events.EMAIL_SEND,
		})

		return response.WriteResponse(ctx, fiber.StatusOK, "Code created successfully.")
	}
}

func RevokeSignupToken(api *types.Api) fiber.Handler {
	logger := utils.Logger()
	db := api.Deps.DB
	res200 := response.NewResponse(fiber.StatusOK, "Token was revoked successfully.")
	return func(ctx *fiber.Ctx) error {
		token := strings.TrimSpace(ctx.FormValue("token"))
		// TODO: Add authorization checks for token creation roles
		if token == "" {
			return response.From(ctx, res200)
		}

		err := db.BackendUsers().DeleteRegToken(ctx.Context(), token)
		if err != nil {
			logger.Error("Error deleting registration token", zap.Error(err))
			return response.WriteResponse(ctx, fiber.StatusInternalServerError, "Something went wrong while processing your request, please try again.")
		}
		return response.From(ctx, res200)
	}
}
