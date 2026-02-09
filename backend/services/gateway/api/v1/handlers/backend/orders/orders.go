package orders

import (
	"errors"
	"strconv"

	"github.com/chibx/vuecom/backend/services/gateway/api/v1/request"
	"github.com/chibx/vuecom/backend/services/gateway/api/v1/response"
	"github.com/chibx/vuecom/backend/services/gateway/internal/types"

	"github.com/chibx/vuecom/backend/shared/models/db/orders"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func CreateOrder(ctx *fiber.Ctx, api *types.Api) error {
	db := api.Deps.DB
	order := orders.Order{}

	err := ctx.BodyParser(&order)

	if err != nil {
		return response.WriteResponse(ctx, fiber.StatusBadRequest, "Invalid request body", nil)
	}

	err = db.Orders().CreateOrder(ctx.Context(), &order)

	if err != nil {
		return response.WriteResponse(ctx, fiber.StatusInternalServerError, "Failed to create order", nil)
	}

	// return ctx.Status(fiber.StatusCreated).SendString("Order Created Succesfully")
	return response.WriteResponse(ctx, fiber.StatusCreated, "Order Created Successfully", nil)
}

func UpdateOrder(ctx *fiber.Ctx, api *types.Api) error {
	return response.WriteResponse(ctx, fiber.StatusOK, "Order updated successfully", nil)
}

func GetOrder(ctx *fiber.Ctx, api *types.Api) error {
	db := api.Deps.DB
	toGet := request.OnlyID{}
	err := ctx.ParamsParser(&toGet)

	if err != nil {
		return response.WriteResponse(ctx, fiber.StatusBadRequest, "Invalid order ID", nil)
	}

	if toGet.ID <= 0 {
		return response.WriteResponse(ctx, fiber.StatusBadRequest, "Order ID cannot be less than 1", nil)
	}

	order, err := db.Orders().GetOrderById(ctx.Context(), toGet.ID)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response.WriteResponse(ctx, fiber.StatusNotFound, "Order with ID "+strconv.Itoa(toGet.ID)+" not found", nil)
		}
		return response.WriteResponse(ctx, fiber.StatusInternalServerError, "Failed to get order", nil)
	}

	return response.WriteResponse(ctx, fiber.StatusOK, "", order)
}

func ListOrders(ctx *fiber.Ctx, api *types.Api) error {
	return response.WriteResponse(ctx, fiber.StatusOK, "", []orders.Order{})
}

func DeleteOrder(ctx *fiber.Ctx, api *types.Api) error {
	return response.WriteResponse(ctx, fiber.StatusOK, "Order deleted successfully", nil)
}

func DeleteOrders(ctx *fiber.Ctx, api *types.Api) error {
	return response.WriteResponse(ctx, fiber.StatusOK, "Orders deleted successfully", nil)
}
