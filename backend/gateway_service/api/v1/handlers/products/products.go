package products

import (
	"vuecom/gateway/api/v1/request"
	"vuecom/gateway/internal/v1/types"

	dbModel "vuecom/shared/models/db"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func CreateProduct(ctx *fiber.Ctx, api *types.Api) error {
	db := api.Deps.DB
	product := dbModel.Product{}

	err := ctx.BodyParser(&product)

	if err != nil {
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}

	err = gorm.G[dbModel.Product](db).Create(ctx.Context(), &product)

	if err != nil {
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}

	return ctx.Status(fiber.StatusCreated).SendString("Product Created Succesfully")
}

func UpdateProduct(ctx *fiber.Ctx) error {
	return nil
}

func GetProduct(ctx *fiber.Ctx, api *types.Api) error {
	db := api.Deps.DB
	toGet := request.OnlyID{}
	// err := ctx.BodyParser(&toGet)
	err := ctx.ParamsParser(&toGet)

	if err != nil {
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}

	if toGet.ID <= 0 {
		return ctx.SendStatus(fiber.StatusBadRequest)
	}

	product, err := gorm.G[dbModel.Product](db).Where("id = ?", toGet.ID).First(ctx.Context())

	if err != nil {
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}

	return ctx.JSON(product)
}

func ListProducts(ctx *fiber.Ctx, api *types.Api) error {
	return nil
}

func DeleteProduct(ctx *fiber.Ctx) error {
	return nil
}

func DeleteProducts(ctx *fiber.Ctx, api *types.Api) error {
	return nil
}
