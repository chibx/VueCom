package inventory

import (
	"github.com/chibx/vuecom/backend/services/gateway/api/v1/response"
	"github.com/chibx/vuecom/backend/services/gateway/internal/global"
	igrpc "github.com/chibx/vuecom/backend/services/gateway/internal/grpc"
	"github.com/chibx/vuecom/backend/services/gateway/internal/types"
	inventoryPr "github.com/chibx/vuecom/backend/shared/proto/go/inventory"
	"go.uber.org/zap"

	"github.com/gofiber/fiber/v2"
)

func HasAnyWarehouse(api *types.Api) fiber.Handler {
	logger := global.Logger
	return func(c *fiber.Ctx) error {
		existsResp, err := igrpc.InventoryClient.HasAnyWarehouse(c.Context(), &inventoryPr.WarehouseExistReq{})
		if err != nil {
			logger.Error("inventory warehouse check failed", zap.Error(err))
			return response.WriteResponse(c, fiber.StatusInternalServerError, "Error occurred while checking inventory warehouses")
		}

		return response.WriteResponse(c, fiber.StatusOK, "", map[string]bool{"exists": existsResp.Exists})
	}
}

func ListWarehouses(api *types.Api) fiber.Handler {
	logger := global.Logger
	return func(c *fiber.Ctx) error {
		listResp, err := igrpc.InventoryClient.ListWarehouses(c.Context(), &inventoryPr.ListWarehousesReq{})
		if err != nil {
			logger.Error("failed to list warehouses", zap.Error(err))
			return response.WriteResponse(c, fiber.StatusInternalServerError, "Error occurred while listing warehouses")
		}

		return response.WriteResponse(c, fiber.StatusOK, "", listResp.Warehouses)
	}
}
