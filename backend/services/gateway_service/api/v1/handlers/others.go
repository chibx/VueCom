package handlers

import (
	"fmt"
	"net/http"
	"slices"
	"strings"
	"vuecom/gateway/internal/types"

	// "github.com/cloudinary/cloudinary-go/v2"
	// "github.com/redis/go-redis/v9"
	// "gorm.io/gorm"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
)

// For validating the admin slug
// func ValidateSlug(ctx *fiber.Ctx, api *types.Api) error {
// 	var routeParts = utils.ExtractRouteParts(ctx.Path())
// 	// var partLen = len(routeParts)
// 	// /"" OR /"admin"
// 	var adminPart = routeParts[0]

// 	if slices.Contains(api.Config.AllowedPaths, adminPart) {
// 		return ctx.Next()
// 	}

// 	// TODO: Remove when true setting up
// 	if adminPart != api.AdminSlug {
// 		ctx.Set(fiber.HeaderContentType, fiber.MIMETextHTMLCharsetUTF8)
// 		return ctx.Status(404).SendString(Page_404)
// 	}

// 	return utils.ServeIndex(ctx)
// }

func serveFile(ctx *fiber.Ctx, path string) error {
	return filesystem.SendFile(ctx, http.Dir("./dist"), path)
}

/**
 * For now our my goal is to focus on the backend entry only, meaning I don't plan on adding validation for customers
 */
func ValidateSlug(ctx *fiber.Ctx, api *types.Api) error {
	adminParam := ctx.Params("admin")

	if adminParam != api.AdminSlug {
		// TODO: Implement admin slug validation logic here
		return fiber.ErrNotFound
	}

	routeParts := strings.Split(strings.TrimRight(ctx.Params("*"), "/"), "/")
	fmt.Println(routeParts)
	if len(routeParts) == 0 {
		err := serveFile(ctx, "index.html")
		if err != nil {
			return fiber.ErrNotFound
		}
		return nil
	}
	var err error
	if len(routeParts) >= 1 {
		fmt.Println("Extra route:", len(slices.Clone(routeParts)[1:]))
		err = serveFile(ctx, strings.Join(routeParts[1:], "/"))
	}
	if err != nil {
		// Serve a catch all to the index.html and handle the issue at the frontend as expected
		err = serveFile(ctx, "index.html")
		if err != nil {
			return fiber.ErrNotFound
		}
	}

	return nil
}
