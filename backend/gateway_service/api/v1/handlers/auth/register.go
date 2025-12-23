package auth

import (
	"vuecom/gateway/internal/types"

	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(app fiber.Router, api *types.Api) {
	auth := app.Group("/auth")
	auth.Post("/register")
	auth.Post("/login")
	auth.Post("/refresh")

	backendAuth := auth.Group("/backend")
	backendAuth.Post("/register")
	backendAuth.Post("/login")
	backendAuth.Post("/refresh")
}
