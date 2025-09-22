package api

import (
	"fmt"
	"slices"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"github.com/gofiber/fiber/v2"
)

type Api struct {
	DB    gorm.DB
	Redis redis.Client
}

// For validating the admin slug
func (s *Api) ValidateSlug(ctx *fiber.Ctx) error {
	var routeParts = extractRouteParts(ctx.Path())
	// var partLen = len(routeParts)
	// /"" OR /"admin"
	var adminPart string = routeParts[0]

	if slices.Contains(AllowedPaths, adminPart) {
		return ctx.Next()
	}

	if adminPart != adminSlug {
		ctx.Context().SetContentType(fiber.MIMETextHTMLCharsetUTF8)
		return ctx.Status(404).SendString(Page_404)
	}

	return serveIndex(ctx)
}

func (a *Api) ApiHandler(ctx *fiber.Ctx) error {
	fmt.Println(ctx.Params("name"))
	return nil
}
