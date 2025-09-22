package api

import "github.com/gofiber/fiber/v2"

func (a *Api) GetProducts(ctx *fiber.Ctx) error {

	return ctx.SendString("Product Yo")
}
