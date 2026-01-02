package orders

import (
	"vuecom/gateway/api/v1/request"
	"vuecom/gateway/internal/types"

	dbModel "vuecom/shared/models/db"

	"github.com/gofiber/fiber/v2"
)

func CreateOrder(ctx *fiber.Ctx, api *types.Api) error {
	db := api.Deps.DB
	order := dbModel.Order{}

	err := ctx.BodyParser(&order)

	if err != nil {
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}

	err = db.Orders().CreateOrder(&order, ctx.Context())

	if err != nil {
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}

	return ctx.Status(fiber.StatusCreated).SendString("Order Created Succesfully")
}

func UpdateOrder(ctx *fiber.Ctx) error {
	return nil
}

func GetOrder(ctx *fiber.Ctx, api *types.Api) error {
	db := api.Deps.DB
	toGet := request.OnlyID{}
	err := ctx.ParamsParser(&toGet)

	if err != nil {
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}

	if toGet.ID <= 0 {
		return ctx.SendStatus(fiber.StatusBadRequest)
	}

	order, err := db.Orders().GetOrderById(toGet.ID, ctx.Context())

	if err != nil {
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}

	return ctx.JSON(order)
}

func ListOrders(ctx *fiber.Ctx, api *types.Api) error {
	return nil
}

func DeleteOrder(ctx *fiber.Ctx) error {
	return nil
}

func DeleteOrders(ctx *fiber.Ctx, api *types.Api) error {
	return nil
}
