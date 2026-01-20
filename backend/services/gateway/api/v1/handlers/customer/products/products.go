package products

import (
	"errors"
	"strconv"

	"github.com/chibx/vuecom/backend/services/gateway/api/v1/request"
	"github.com/chibx/vuecom/backend/services/gateway/api/v1/response"
	"github.com/chibx/vuecom/backend/services/gateway/internal/types"

	catModels "github.com/chibx/vuecom/backend/shared/models/db/catalog"

	"github.com/gofiber/fiber/v2"
)

func CreateProduct(ctx *fiber.Ctx, api *types.Api) error {
	db := api.Deps.DB
	product := catModels.Product{}

	err := ctx.BodyParser(&product)

	if err != nil {
		return response.NewResponse(ctx, fiber.StatusBadRequest, "Validation error")
	}

	err = db.Products().CreateProduct(&product, ctx.Context())

	if err != nil {
		return response.NewResponse(ctx, fiber.StatusInternalServerError, "Error occurred creating product")
	}

	return response.NewResponse(ctx, fiber.StatusCreated, "Product Created Succesfully")
}

func UpdateProduct(ctx *fiber.Ctx) error {
	return response.NewResponse(ctx, fiber.StatusOK, "", nil)
}

func GetProduct(ctx *fiber.Ctx, api *types.Api) error {
	db := api.Deps.DB
	toGet := request.OnlyID{}
	// err := ctx.BodyParser(&toGet)
	err := ctx.ParamsParser(&toGet)

	if err != nil {
		return response.NewResponse(ctx, fiber.StatusBadRequest, "Validation error")
	}

	if toGet.ID <= 0 {
		return response.NewResponse(ctx, fiber.StatusBadRequest, "Product ID cannot be less than 1")
	}

	// product, err := gorm.G[dbModel.Product](db).Where("id = ?", toGet.ID).First(ctx.Context())
	product, err := db.Products().GetProductById(toGet.ID, ctx.Context())

	if err != nil {
		if errors.Is(err, types.ErrDbNil) {
			return response.NewResponse(ctx, fiber.StatusNotFound, "Product with ID "+strconv.Itoa(toGet.ID)+" not found")
		}
		return response.NewResponse(ctx, fiber.StatusInternalServerError, "Error occurred while fetching product")
	}

	return response.NewResponse(ctx, fiber.StatusOK, "", product)
}

func ListProducts(ctx *fiber.Ctx, api *types.Api) error {

	return response.NewResponse(ctx, fiber.StatusOK, "", nil)
}

func DeleteProduct(ctx *fiber.Ctx) error {

	return response.NewResponse(ctx, fiber.StatusOK, "", nil)
}

func DeleteProducts(ctx *fiber.Ctx, api *types.Api) error {
	return response.NewResponse(ctx, fiber.StatusOK, "", nil)
}
