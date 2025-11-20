package v1

import (
	"vuecom/shared/models"

	dbModel "vuecom/shared/models/db"

	"github.com/gofiber/fiber/v2"
)

func (api *Api) CreateProduct(ctx *fiber.Ctx) error {
	product := dbModel.Product{}

	err := ctx.BodyParser(&product)

	if err != nil {
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}

	result := api.DB.Create(&product)

	if result.Error != nil {
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}

	return ctx.Status(fiber.StatusCreated).SendString("Product Created Succesfully")
}

func UpdateProduct(ctx *fiber.Ctx) error {
	return nil
}

func (api *Api) GetProduct(ctx *fiber.Ctx) error {
	toGet := models.OnlyID{}
	// err := ctx.BodyParser(&toGet)
	err := ctx.ParamsParser(&toGet)

	if err != nil {
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}

	if toGet.ID <= 0 {
		return ctx.SendStatus(fiber.StatusBadRequest)
	}

	product := dbModel.Product{}
	result := api.DB.Model(dbModel.Product{}).Where("id = ?", toGet.ID).First(&product)

	if result.Error != nil {
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}

	return ctx.JSON(product)
}

func (api *Api) ListProducts(ctx *fiber.Ctx) error {
	return nil
}

func (api *Api) DeleteProduct(ctx *fiber.Ctx) error {
	return nil
}

func (api *Api) DeleteProducts(ctx *fiber.Ctx) error {
	return nil
}
