package api

import (
	"vuecom/models"
	dbModel "vuecom/models/db"

	"github.com/gofiber/fiber/v2"
)

func CreateProduct(ctx *fiber.Ctx) error {
	return nil
}

func UpdateProduct(ctx *fiber.Ctx) error {
	return nil
}

func (api *Api) GetProducts(ctx *fiber.Ctx) error {
	toGet := models.OnlyID{}
	err := ctx.BodyParser(&toGet)

	if err != nil {
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}

	product := dbModel.Product{}
	result := api.DB.First(&product)

	return ctx.SendString("Product Yo")
}

func ListProducts(ctx *fiber.Ctx) error {
	return nil
}

func DeleteProduct(ctx *fiber.Ctx) error {
	return nil
}

func DeleteProducts(ctx *fiber.Ctx) error {
	return nil
}
