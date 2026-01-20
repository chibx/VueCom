package customer

import (
	"github.com/chibx/vuecom/backend/services/gateway/api/v1/middlewares"
	"github.com/chibx/vuecom/backend/services/gateway/internal/types"

	"github.com/gofiber/fiber/v2"

	orderHandler "github.com/chibx/vuecom/backend/services/gateway/api/v1/handlers/customer/orders"
	productHandler "github.com/chibx/vuecom/backend/services/gateway/api/v1/handlers/customer/products"
)

func LoadRoutes(app fiber.Router, api *types.Api) {
	// app.Use(middlewares.AuthMiddleware(api))

	/* /v1 handlers */
	v1 := app.Group("/api/customer", middlewares.CustomerRateLimit(api))
	productHandler.RegisterRoutes(v1, api)
	orderHandler.RegisterRoutes(v1, api)
}
