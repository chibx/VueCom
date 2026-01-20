package backend

import (
	"github.com/chibx/vuecom/backend/services/gateway/api/v1/middlewares"
	"github.com/chibx/vuecom/backend/services/gateway/internal/types"

	"github.com/gofiber/fiber/v2"

	adminHandler "github.com/chibx/vuecom/backend/services/gateway/api/v1/handlers/backend/admin"
	orderHandler "github.com/chibx/vuecom/backend/services/gateway/api/v1/handlers/backend/orders"
	productHandler "github.com/chibx/vuecom/backend/services/gateway/api/v1/handlers/backend/products"
)

func LoadRoutes(app fiber.Router, api *types.Api) {
	// app.Use(middlewares.AuthMiddleware(api))

	/* /v1 handlers */
	v1 := app.Group("/api/backend", middlewares.BackendRateLimit(api))
	productHandler.RegisterRoutes(v1, api)
	adminHandler.RegisterRoutes(v1, api)
	orderHandler.RegisterRoutes(v1, api)
}
