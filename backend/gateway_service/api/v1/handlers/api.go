package handlers

import (
	"slices"
	"vuecom/gateway/internal/v1/types"
	"vuecom/gateway/utils"

	// "github.com/cloudinary/cloudinary-go/v2"
	// "github.com/redis/go-redis/v9"
	// "gorm.io/gorm"
	"github.com/gofiber/fiber/v2"
)

// For validating the admin slug
func ValidateSlug(ctx *fiber.Ctx, api *types.Api) error {
	var routeParts = utils.ExtractRouteParts(ctx.Path())
	// var partLen = len(routeParts)
	// /"" OR /"admin"
	var adminPart string = routeParts[0]

	if slices.Contains(api.Config.AllowedPaths, adminPart) {
		return ctx.Next()
	}

	// TODO: Remove when true setting up
	if adminPart != api.Config.MockAdminSlug {
		ctx.Set(fiber.HeaderContentType, fiber.MIMETextHTMLCharsetUTF8)
		return ctx.Status(404).SendString(Page_404)
	}

	return utils.ServeIndex(ctx)
}
