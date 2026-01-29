package auth

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"net"
	"time"

	serverErrors "github.com/chibx/vuecom/backend/shared/errors/server"
	userModels "github.com/chibx/vuecom/backend/shared/models/db/users"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/chibx/vuecom/backend/services/gateway/internal/constants"
	"github.com/chibx/vuecom/backend/services/gateway/internal/types"
	"github.com/chibx/vuecom/backend/services/gateway/internal/utils"

	"github.com/gofiber/fiber/v2"
)

// GenerateRefreshToken (as before).
func GenerateRefreshToken() (string, error) {
	// ... (crypto/rand + hex)
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func CompositeRefreshToken() (token string, refreshHash string, err error) {
	refreshToken, err := GenerateRefreshToken()
	if err != nil {
		return "", "", err
	}
	refreshTokenHash, err := GenerateHashFromString(refreshToken, DefaultHashParams)
	if err != nil {
		return "", "", err
	}

	return fmt.Sprintf("%s", refreshToken), refreshTokenHash, nil
}

func ValidateBackendUserSess(ctx *fiber.Ctx, session *userModels.BackendSession) error {
	created_at := session.CreatedAt
	lastIp := session.LastIP
	// Handle user_agent validation later
	_ = session.UserAgent

	// TODO: Add validation logic for expiry and IP address
	_ = created_at
	_ = lastIp

	current_time := time.Now()
	if current_time.Sub(created_at) > constants.BackendRefreshTkDur {
		return serverErrors.NewSessionErr(serverErrors.SessionExpired, "Session has expired")
	}
	// CAUTION: This is a basic IP check.
	ip := net.ParseIP(ctx.IP())
	if ip == nil {
		return serverErrors.NewSessionErr(serverErrors.SessionInvalidIpAddr, "IP address is either missing or invalid")
	}

	if ip.String() != lastIp {
		// TODO: Maybe check if the IP range is valid and secure instead of rejecting it outright
		return serverErrors.NewSessionErr(serverErrors.SessionDiffIpAddr, "IP address does not match")
	}

	return nil
}

func CreateBackendSession(ctx context.Context, session *userModels.BackendSession, api *types.Api) error {
	db := api.Deps.DB
	var err error
	logger := utils.Logger()
	err = db.BackendUsers().CreateSession(ctx, session)
	if err != nil {
		logger.Error("Failed to create backend session", zap.Error(err))
	}

	return err
}

func GetBackendUserSession(ctx context.Context, tokenId string, api *types.Api) (*userModels.BackendSession, error) {
	db := api.Deps.DB
	logger := utils.Logger()

	var backend_session *userModels.BackendSession

	backend_session, err := db.BackendUsers().GetSessionByTokenId(ctx, tokenId)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Error("backend user session not found in db", zap.Error(err))
			return nil, serverErrors.NewServerErr(fiber.StatusUnauthorized, "User Session not found. Consider logging in again")
		}
		logger.Error("failed to get backend user session from db", zap.Error(err))
		return nil, serverErrors.NewServerErr(fiber.StatusInternalServerError, "Something went wrong while fetching your session data. Please try again later.")
	}

	return backend_session, nil
}

func GetCustomerSession(ctx context.Context, tokenId string, api *types.Api) (*userModels.CustomerSession, error) {
	db := api.Deps.DB
	logger := utils.Logger()

	var backend_session *userModels.CustomerSession

	backend_session, err := db.Customers().GetSessionByTokenId(ctx, tokenId)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Error("backend user session not found in db", zap.Error(err))
			return nil, serverErrors.NewServerErr(fiber.StatusUnauthorized, "User Session not found. Consider logging in again")
		}
		logger.Error("failed to get backend user session from db", zap.Error(err))
		return nil, serverErrors.NewServerErr(fiber.StatusInternalServerError, "Something went wrong while fetching your session data. Please try again later.")
	}

	return backend_session, nil
}
