package auth

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/chibx/vuecom/backend/services/gateway/internal/auth"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// Setup DB and Redis.
var db *gorm.DB
var rdb *redis.Client

// RefreshToken model (as before).
type RefreshToken struct {
	gorm.Model
	Token     string `gorm:"uniqueIndex"`
	UserID    uint
	ExpiresAt time.Time
}

// Secrets.
var accessSecret = []byte("your-access-secret-key")

// Login handler: Issues both tokens.
func BackendLoginJWT(ctx *fiber.Ctx) error {
	// Fake auth: Assume userID from credentials.
	userID := uint(123)

	// Access token.
	accessClaims := jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(15 * time.Minute).Unix(),
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessSigned, err := accessToken.SignedString(accessSecret)
	if err != nil {
		return err
	}

	// Refresh token.
	refresh, err := auth.GenerateRefreshToken()
	if err != nil {
		return err
	}
	expiry := time.Now().Add(7 * 24 * time.Hour)
	if err := db.Create(&RefreshToken{Token: refresh, UserID: userID, ExpiresAt: expiry}).Error; err != nil {
		return err
	}

	if err := rdb.Set(ctx.Context(), "refresh:"+refresh, userID, time.Until(expiry)).Err(); err != nil {
		return err
	}

	// Set refresh cookie.
	ctx.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    refresh,
		Expires:  expiry,
		HTTPOnly: true,
		Secure:   true,
		SameSite: "Strict",
	})

	return ctx.JSON(map[string]string{"access_token": accessSigned})
}

// Refresh handler: Uses refresh to create new access.
func RefreshJWT(ctx *fiber.Ctx) error {
	refresh := ctx.Cookies("refresh_token")
	if refresh == "" {
		return ctx.Status(http.StatusUnauthorized).SendString("No refresh token")
	}

	var userIDStr string
	var refToken RefreshToken

	// Step 1: Check Redis (cache) for fast validation.
	userIDStr, err := rdb.Get(ctx.Context(), "refresh:"+refresh).Result()
	if !errors.Is(err, redis.Nil) {
		// Cache miss: Fall back to DB.
		if err := db.Where("token = ? AND expires_at > ?", refresh, time.Now()).First(&refToken).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ctx.Status(http.StatusUnauthorized).SendString("Invalid or expired refresh token")
			}
			return err
		}
		// Repopulate Redis.
		if err := rdb.Set(ctx.Context(), "refresh:"+refresh, refToken.UserID, time.Until(refToken.ExpiresAt)).Err(); err != nil {
			return err
		}
		userIDStr = fmt.Sprintf("%d", refToken.UserID)
	} else if err != nil {
		return err
	}

	// Step 2: Parse userID.
	userID := uint(0) // Convert string to uint.
	if _, err := fmt.Sscanf(userIDStr, "%d", &userID); err != nil {
		return err
	}

	// Step 3: Generate new short-lived access token using the validated userID.
	accessClaims := jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(15 * time.Minute).Unix(), // Short-lived.
		// You can add fresh claims here, e.g., updated roles from DB.
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessSigned, err := accessToken.SignedString(accessSecret)
	if err != nil {
		return err
	}

	// Optional: Rotate refresh for extra security (generate new refresh, delete old).
	// But for simplicity, we reuse the existing one.

	return ctx.JSON(map[string]string{"access_token": accessSigned})
}

// Logout: Revoke refresh.
func LogoutJWT(ctx *fiber.Ctx) error {
	refresh := ctx.Cookies("refresh_token")
	if refresh != "" {
		rdb.Del(ctx.Context(), "refresh:"+refresh)
		db.Where("token = ?", refresh).Delete(&RefreshToken{})
	}
	ctx.ClearCookie("refresh_token")
	return ctx.SendString("Logged out")
}
