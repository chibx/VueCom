package auth

import (
	"crypto/rand"
	"encoding/hex"
	"time"
	"vuecom/gateway/internal/constants"
	serverErr "vuecom/shared/errors/server"
	dbModels "vuecom/shared/models/db"

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

func ValidateBackendUserSess(ctx *fiber.Ctx, session *dbModels.BackendSession) error {
	created_at := session.CreatedAt
	ipAddr := session.IpAddr
	// Handle user_agent validation later
	_ = session.UserAgent

	// TODO: Add validation logic for expiry and IP address
	_ = created_at
	_ = ipAddr

	current_time := time.Now()
	if current_time.Sub(created_at) > constants.BACKEND_SESSION_TIMEOUT {
		return serverErr.NewSessionErr(serverErr.SessionExpired, "Session has expired")
	}
	// CAUTION: This is a basic IP check.
	if ctx.IP() != ipAddr {
		return serverErr.NewSessionErr(serverErr.SessionDiffIpAddr, "IP address does not match")
	}

	return nil
}
