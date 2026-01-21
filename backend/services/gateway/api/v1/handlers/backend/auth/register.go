package auth

import (
	"github.com/chibx/vuecom/backend/services/gateway/internal/types"

	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(app fiber.Router, api *types.Api) {
	auth := app.Group("/auth")

	auth.Post("/register", Register(api))
	auth.Post("/login", Login(api))
}
