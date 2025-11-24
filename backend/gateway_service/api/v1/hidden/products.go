package v1

// import (
// 	"vuecom/shared/models"

// 	dbModel "vuecom/shared/models/db"

// 	"github.com/gofiber/fiber/v2"
// 	"gorm.io/gorm"
// )

// func (api *Api) CreateProduct(ctx *fiber.Ctx) error {
// 	db := api.Deps.DB
// 	product := dbModel.Product{}

// 	err := ctx.BodyParser(&product)

// 	if err != nil {
// 		return ctx.SendStatus(fiber.StatusInternalServerError)
// 	}

// 	err = gorm.G[dbModel.Product](db).Create(ctx.Context(), &product)

// 	if err != nil {
// 		return ctx.SendStatus(fiber.StatusInternalServerError)
// 	}

// 	return ctx.Status(fiber.StatusCreated).SendString("Product Created Succesfully")
// }

// func UpdateProduct(ctx *fiber.Ctx) error {
// 	return nil
// }

// func (api *Api) GetProduct(ctx *fiber.Ctx) error {
// 	db := api.Deps.DB
// 	toGet := models.OnlyID{}
// 	// err := ctx.BodyParser(&toGet)
// 	err := ctx.ParamsParser(&toGet)

// 	if err != nil {
// 		return ctx.SendStatus(fiber.StatusInternalServerError)
// 	}

// 	if toGet.ID <= 0 {
// 		return ctx.SendStatus(fiber.StatusBadRequest)
// 	}

// 	product, err := gorm.G[dbModel.Product](db).Where("id = ?", toGet.ID).First(ctx.Context())

// 	if err != nil {
// 		return ctx.SendStatus(fiber.StatusInternalServerError)
// 	}

// 	return ctx.JSON(product)
// }

// func (api *Api) ListProducts(ctx *fiber.Ctx) error {
// 	return nil
// }

// func (api *Api) DeleteProduct(ctx *fiber.Ctx) error {
// 	return nil
// }

// func (api *Api) DeleteProducts(ctx *fiber.Ctx) error {
// 	return nil
// }
