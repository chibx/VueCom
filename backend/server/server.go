package server

import (
	"os"
	"slices"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

var adminSlug = "admin123"

type Server struct {
	DB    gorm.DB
	Redis redis.Client
}

func WriteFile(ctx *fiber.Ctx, path string, ctype string) error {
	file, err := os.ReadFile(path)

	if err != nil {
		log.Error(err)
		return fiber.ErrNotFound
	}

	header := &ctx.Response().Header
	header.SetContentType(ctype)
	header.SetContentLength(len(file))

	_, err = ctx.Write(file)

	if err != nil {
		log.Error(err)
		return fiber.ErrInternalServerError
	}

	return nil
}

func (s *Server) ValidateSlug(ctx *fiber.Ctx) error {
	var routeParts = extractRouteParts(ctx.Path())
	// var partLen = len(routeParts)
	// /"" OR /"admin"
	var adminPart string = routeParts[0]

	if slices.Contains(AllowedPaths, adminPart) {
		return ctx.Next()
	}

	if adminPart != adminSlug {
		return ctx.SendStatus(fiber.StatusNotFound)
	}

	return s.serveIndex(ctx)
}

func (s *Server) serveIndex(ctx *fiber.Ctx) error {
	return WriteFile(ctx, "./dist/index.html", fiber.MIMETextHTMLCharsetUTF8)
}
