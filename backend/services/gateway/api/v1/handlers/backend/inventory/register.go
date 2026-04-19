package inventory

import (
	"github.com/chibx/vuecom/backend/services/gateway/internal/types"

	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(app fiber.Router, api *types.Api) {
	inventoryGroup := app.Group("/inventory")

	inventoryGroup.Get("/warehouses/exist", HasAnyWarehouse(api))
	inventoryGroup.Get("/warehouses", ListWarehouses(api))
	inventoryGroup.Post("/warehouses", CreateWarehouse(api))
	inventoryGroup.Delete("/warehouses", DeleteWarehouse(api))
	inventoryGroup.Post("/stock-movements", CreateStockMovement(api))
	inventoryGroup.Get("/stock-movements", ListStockMovements(api))
}
