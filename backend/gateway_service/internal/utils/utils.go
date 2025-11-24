package utils

import (
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"

	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/utils"
)

func ExtractRouteParts(route string) []string {
	route = utils.CopyString(route)
	var routeLength = len(route)
	var routeParts []string

	if routeLength == 1 {
		// It's just "/"
		routeParts = []string{""}
	} else {
		var hasTrailingSlash = string(route[routeLength-1]) == "/"

		if hasTrailingSlash {
			route = route[:len(route)-1]
		}
		// First index is gonna be ""
		routeParts = strings.Split(route, "/")[1:]
	}

	return routeParts
}

func ServeIndex(ctx *fiber.Ctx) error {
	return writeFile(ctx, "./dist/index.html", fiber.MIMETextHTMLCharsetUTF8)
}

func writeFile(ctx *fiber.Ctx, path string, ctype string) error {
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
