package api

import (
	"fmt"
	"slices"
	"vuecom/api/utils"
	"vuecom/config"
	"vuecom/models"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"github.com/gofiber/fiber/v2"
)

type Api struct {
	DB    *gorm.DB
	Redis redis.Client
}

// For validating the admin slug
func (s *Api) ValidateSlug(ctx *fiber.Ctx) error {
	var routeParts = utils.ExtractRouteParts(ctx.Path())
	// var partLen = len(routeParts)
	// /"" OR /"admin"
	var adminPart string = routeParts[0]

	if slices.Contains(config.AllowedPaths, adminPart) {
		return ctx.Next()
	}

	// TODO: Remove when true setting up
	if adminPart != config.MockAdminSlug {
		ctx.Set(fiber.HeaderContentType, fiber.MIMETextHTMLCharsetUTF8)
		return ctx.Status(404).SendString(Page_404)
	}

	return utils.ServeIndex(ctx)
}

func (a *Api) ApiHandler(ctx *fiber.Ctx) error {
	fmt.Println(ctx.Params("name"))
	return nil
}

func (api *Api) initializeApp(ctx *fiber.Ctx) error {
	appData := models.CreateAppData{}
	err := ctx.BodyParser(&appData)

	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	return nil
}
