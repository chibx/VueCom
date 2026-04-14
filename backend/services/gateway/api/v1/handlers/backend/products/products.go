package products

import (
	"errors"
	"strconv"

	"github.com/chibx/vuecom/backend/services/gateway/api/v1/request"
	"github.com/chibx/vuecom/backend/services/gateway/api/v1/response"
	"github.com/chibx/vuecom/backend/services/gateway/internal/global"
	igrpc "github.com/chibx/vuecom/backend/services/gateway/internal/grpc"
	"github.com/chibx/vuecom/backend/services/gateway/internal/types"
	"github.com/chibx/vuecom/backend/services/gateway/internal/utils"
	"go.uber.org/zap"

	reqTypes "github.com/chibx/vuecom/backend/services/gateway/api/v1/request/catalog"

	serverErrors "github.com/chibx/vuecom/backend/shared/errors/server"
	"github.com/gofiber/fiber/v2"
)

func CreateProduct(api *types.Api) fiber.Handler {
	err500 := fiber.NewError(fiber.StatusInternalServerError, "Error occurred while creating product, please try again.")
	logger := global.Logger()
	return func(c *fiber.Ctx) error {
		var err error

		reqBody := reqTypes.CreateProductReq{}

		err = c.BodyParser(&reqBody)
		if err != nil {
			return response.WriteResponse(c, fiber.StatusBadRequest, "Validation error")
		}

		err = utils.Validator().Struct(reqBody)
		isFatal, errorBag := serverErrors.HandleValidationError(err)
		if isFatal {
			logger.Error("InvalidValidationError while creating a signup token", zap.Error(err))
			return response.WriteResponse(c, fiber.ErrBadRequest.Code, err500.Message)
		}
		if len(errorBag) > 0 {
			return response.WriteResponse(c, fiber.StatusBadRequest, "One or more fields are invalid", errorBag)
		}

		normalizeProdReq(&reqBody)

		prodRpc, err := createProductToRpc(&reqBody)
		if err != nil {
			return response.FromFiberError(c, err500)
		}
		prodRpcResp, err := igrpc.CatalogClient.CreateProduct(c.Context(), prodRpc)
		_ = prodRpcResp.Id

		if err != nil {
			return response.WriteResponse(c, fiber.StatusInternalServerError, "Error occurred creating product")
		}

		return response.WriteResponse(c, fiber.StatusCreated, "Product Created Succesfully")
	}
}

func UpdateProduct(ctx *fiber.Ctx) error {
	return response.WriteResponse(ctx, fiber.StatusOK, "", nil)
}

func GetProduct(ctx *fiber.Ctx, api *types.Api) error {
	db := api.Deps.DB
	toGet := request.OnlyID{}
	// err := ctx.BodyParser(&toGet)
	err := ctx.ParamsParser(&toGet)

	if err != nil {
		return response.WriteResponse(ctx, fiber.StatusBadRequest, "Validation error")
	}

	if toGet.ID <= 0 {
		return response.WriteResponse(ctx, fiber.StatusBadRequest, "Product ID cannot be less than 1")
	}

	// product, err := gorm.G[dbModel.Product](db).Where("id = ?", toGet.ID).First(ctx.Context())
	product, err := db.Products().GetProductById(ctx.Context(), toGet.ID)

	if err != nil {
		if errors.Is(err, serverErrors.ErrDBRecordNotFound) {
			return response.WriteResponse(ctx, fiber.StatusNotFound, "Product with ID "+strconv.Itoa(toGet.ID)+" not found")
		}
		return response.WriteResponse(ctx, fiber.StatusInternalServerError, "Error occurred while fetching product")
	}

	return response.WriteResponse(ctx, fiber.StatusOK, "", product)
}

func ListProducts(ctx *fiber.Ctx, api *types.Api) error {

	return response.WriteResponse(ctx, fiber.StatusOK, "", nil)
}

func DeleteProduct(ctx *fiber.Ctx) error {

	return response.WriteResponse(ctx, fiber.StatusOK, "", nil)
}

func DeleteProducts(ctx *fiber.Ctx, api *types.Api) error {
	return response.WriteResponse(ctx, fiber.StatusOK, "", nil)
}
