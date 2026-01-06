package constants

import "time"

// Simple way of doing this, a better way would be to use a manifest from a bundler
var PublicAssets = []string{
	"styles.css",
	"assets",
	"robots.txt",
}

const BACKEND_SESSION_TIMEOUT = 30 * time.Minute
const PublicFolder = "dist"
const ApiKeyCtxKey = "api_key"
const BackendUserCtxKey = "backend_user"
const BackendCookieKey = "backend_session"
