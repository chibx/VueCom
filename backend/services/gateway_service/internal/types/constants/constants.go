package constants

import (
	"os"
	"time"
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

const BACKEND_SESSION_TIMEOUT = 30 * time.Minute
const ApiKeyCtxKey = "api_key"
const BackendUserCtxKey = "backend_user"
const BackendCookieKey = "backend_session"
