package inventory

import (
	invReq "github.com/chibx/vuecom/backend/services/gateway/api/v1/request/inventory"
	"github.com/chibx/vuecom/backend/services/gateway/api/v1/response"
	invResp "github.com/chibx/vuecom/backend/services/gateway/api/v1/response/inventory"
	"github.com/chibx/vuecom/backend/services/gateway/internal/global"
	igrpc "github.com/chibx/vuecom/backend/services/gateway/internal/grpc"
	"github.com/chibx/vuecom/backend/services/gateway/internal/types"
	"github.com/chibx/vuecom/backend/services/gateway/internal/utils"
	serverErrors "github.com/chibx/vuecom/backend/shared/errors/server"
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
			return response.WriteResponse(c, fiber.StatusInternalServerError, "Error occurred listing warehouses")
		}

		return response.WriteResponse(c, fiber.StatusOK, "", listResp.Warehouses)
	}
}

// func AddProduct(api *types.Api) fiber.Handler {
// 	logger := global.Logger
// 	return func(c *fiber.Ctx) error {
// 		var req inventoryPr.AddProductRequest
// 		err := c.BodyParser(&req)
// 		if err != nil {
// 			return response.WriteResponse(c, fiber.StatusBadRequest, "Invalid request body")
// 		}

// 		resp, err := igrpc.InventoryClient.CreateProductRecord(c.Context(), &req)
// 		if err != nil {
// 			logger.Error("failed to create product record", zap.Error(err))
// 			return response.WriteResponse(c, fiber.StatusInternalServerError, "Error occurred while creating product record")
// 		}

// 		return response.WriteResponse(c, fiber.StatusCreated, "Product record created", resp)
// 	}
// }

func CreateWarehouse(api *types.Api) fiber.Handler {
	logger := global.Logger
	err500 := fiber.NewError(fiber.StatusInternalServerError, "Error occurred while creating warehouse, please try again.")
	return func(c *fiber.Ctx) error {
		var req invReq.CreateWarehouseReq
		err := c.BodyParser(&req)
		if err != nil {
			return response.WriteResponse(c, fiber.StatusBadRequest, "Invalid request body")
		}

		err = utils.Validator().Struct(req)
		isFatal, errorBag := serverErrors.HandleValidationError(err)
		if isFatal {
			logger.Error("InvalidValidationError while creating warehouse", zap.Error(err))
			return response.WriteResponse(c, fiber.ErrInternalServerError.Code, err500.Message)
		}
		if len(errorBag) > 0 {
			return response.WriteResponse(c, fiber.StatusBadRequest, "One or more fields are invalid", errorBag)
		}

		protoReq := &inventoryPr.CreateWarehouseReq{
			Code:      req.Code,
			Name:      req.Name,
			Address:   req.Address,
			City:      req.City,
			StateId:   req.StateID,
			CountryId: req.CountryID,
			IsActive:  req.IsActive,
		}

		resp, err := igrpc.InventoryClient.CreateWarehouse(c.Context(), protoReq)
		if err != nil {
			logger.Error("failed to create warehouse", zap.Error(err))
			return response.WriteResponse(c, fiber.StatusInternalServerError, "Error occurred while creating warehouse")
		}

		return response.WriteResponse(c, fiber.StatusCreated, "Warehouse created", invResp.CreateWarehouseResp{ID: resp.Id})
	}
}

func DeleteWarehouse(api *types.Api) fiber.Handler {
	logger := global.Logger
	err500 := fiber.NewError(fiber.StatusInternalServerError, "Error occurred while deleting warehouse, please try again.")
	return func(c *fiber.Ctx) error {
		var req invReq.DeleteWarehouseReq
		err := c.BodyParser(&req)
		if err != nil {
			return response.WriteResponse(c, fiber.StatusBadRequest, "Invalid request body")
		}
		warehouseIds := req.WarehouseIDs
		if len(warehouseIds) == 0 {
			return response.WriteResponse(c, fiber.StatusOK, "Nothing to delete")
		}

		protoReq := &inventoryPr.DeleteWarehouseReq{WarehouseIds: warehouseIds}

		_, err = igrpc.InventoryClient.DeleteWarehouse(c.Context(), protoReq)
		if err != nil {
			logger.Error("failed to delete warehouse", zap.Error(err))
			return response.FromFiberError(c, err500)
		}

		return response.WriteResponse(c, fiber.StatusOK, "Warehouse(s) deleted")
	}
}

func CreateStockMovement(api *types.Api) fiber.Handler {
	logger := global.Logger
	err500 := fiber.NewError(fiber.StatusInternalServerError, "Error applying stock movement, please try again.")
	return func(c *fiber.Ctx) error {
		var req invReq.CreateStockMovementReq
		err := c.BodyParser(&req)
		if err != nil {
			return response.WriteResponse(c, fiber.StatusBadRequest, "Invalid request body")
		}

		err = utils.Validator().Struct(req)
		isFatal, errorBag := serverErrors.HandleValidationError(err)
		if isFatal {
			logger.Error("InvalidValidationError, applying stock movement", zap.Error(err))
			return response.WriteResponse(c, fiber.ErrInternalServerError.Code, err500.Message)
		}
		if len(errorBag) > 0 {
			return response.WriteResponse(c, fiber.StatusBadRequest, "One or more fields are invalid", errorBag)
		}

		protoReq := &inventoryPr.CreateStockMovementReq{
			InventoryId:  req.InventoryID,
			Sku:          req.SKU,
			WarehouseId:  req.WarehouseID,
			MovementType: req.MovementType,
			Quantity:     req.Quantity,
			Reference:    req.Reference,
			Notes:        req.Notes,
			CreatedBy:    req.CreatedBy,
		}

		resp, err := igrpc.InventoryClient.CreateStockMovement(c.Context(), protoReq)
		if err != nil {
			logger.Error("failed to create stock movement", zap.Error(err))
			return response.WriteResponse(c, fiber.StatusInternalServerError, "Error occurred while creating stock movement")
		}

		return response.WriteResponse(c, fiber.StatusCreated, "Stock movement created", invResp.CreateStockMovementResp{ID: resp.Id})
	}
}

func ListStockMovements(api *types.Api) fiber.Handler {
	logger := global.Logger
	err500 := fiber.NewError(fiber.StatusInternalServerError, "Error listing stock movements, please try again.")
	return func(c *fiber.Ctx) error {
		var req invReq.ListStockMovementsReq
		err := c.BodyParser(&req)
		if err != nil {
			return response.WriteResponse(c, fiber.StatusBadRequest, "Invalid request body")
		}

		err = utils.Validator().Struct(req)
		isFatal, errorBag := serverErrors.HandleValidationError(err)
		if isFatal {
			logger.Error("InvalidValidationError, applying stock movement", zap.Error(err))
			return response.WriteResponse(c, fiber.ErrInternalServerError.Code, err500.Message)
		}
		if len(errorBag) > 0 {
			return response.WriteResponse(c, fiber.StatusBadRequest, "One or more fields are invalid", errorBag)
		}

		protoReq := &inventoryPr.ListStockMovementsReq{
			WarehouseId: req.WarehouseID,
			Sku:         req.SKU,
		}

		resp, err := igrpc.InventoryClient.ListStockMovements(c.Context(), protoReq)
		if err != nil {
			logger.Error("failed to list stock movements", zap.Error(err))
			return response.WriteResponse(c, fiber.StatusInternalServerError, "Error occurred while listing stock movements")
		}

		var movements []invResp.StockMovement
		for _, m := range resp.StockMovements {
			movements = append(movements, invResp.StockMovement{
				ID:           m.Id,
				InventoryID:  m.InventoryId,
				SKU:          m.Sku,
				WarehouseID:  m.WarehouseId,
				MovementType: m.MovementType,
				Quantity:     m.Quantity,
				Reference:    m.Reference,
				Notes:        m.Notes,
				CreatedBy:    m.CreatedBy,
				CreatedAt:    m.CreatedAt,
			})
		}

		return response.WriteResponse(c, fiber.StatusOK, "", invResp.ListStockMovementsResp{StockMovements: movements})
	}
}
