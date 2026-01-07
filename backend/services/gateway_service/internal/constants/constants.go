package constants

import (
	"os"
	"time"

	"github.com/go-redis/redis_rate/v10"
)

// Simple way of doing this, a better way would be to use a manifest from a bundler
var PublicAssets = []string{
	"styles.css",
	"assets",
	"favicon.ico",
	"robots.txt",
}

var PublicFolder = func() string {
	if folder := os.Getenv("PUBLIC_FOLDER"); folder != "" {
		return folder
	}
	return "dist"
}()

var (
	// Token bucket configurations
	GlobalLimit = redis_rate.Limit{
		Rate:   50000, // Total requests across the entire app
		Period: time.Minute,
		Burst:  100000, // Allow large bursts (e.g., traffic spikes)
	}

	CustomerLimit = redis_rate.Limit{
		Rate:   100, // Per customer
		Period: time.Minute,
		Burst:  200,
	}

	BackendLimit = redis_rate.Limit{
		Rate:   1000, // Higher for devs/admins
		Period: time.Minute,
		Burst:  2000,
	}
)

// Max allowed image size in bytes i.e 5MB
const MAX_IMAGE_UPLOAD = 5 * 1024 * 1024
const GlobalLimitKey = "rl_global:app"
const AnonymousLimitKey = "rl_anonymous:" // With Ip then
const CustomerLimitKey = "rl_customer:"
const CustomerHeaderKey = "X-Customer-Id"
const BackendLimitKey = "rl_backend:"
const BackendSessionTimeout = 30 * time.Minute
const ApiKeyCtxKey = "api_key"
const BackendUserCtxKey = "backend_user"
const BackendCookieKey = "backend_session"
