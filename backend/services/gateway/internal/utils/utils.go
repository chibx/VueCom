package utils

import (
	"errors"
	"io"
	"os"
	"slices"
	"strings"

	"github.com/chibx/vuecom/backend/services/gateway/api/v1/request"

	"github.com/gabriel-vasile/mimetype"
	"github.com/gofiber/fiber/v2"

	"github.com/gofiber/fiber/v2/log"
)

// func ExtractRouteParts(route string) []string {
// 	route = utils.CopyString(route)
// 	var routeLength = len(route)
// 	var routeParts []string

// 	if routeLength == 1 {
// 		// It's just "/"
// 		routeParts = []string{""}
// 	} else {
// 		var hasTrailingSlash = string(route[routeLength-1]) == "/"

// 		if hasTrailingSlash {
// 			route = route[:len(route)-1]
// 		}
// 		// First index is gonna be ""
// 		routeParts = strings.Split(route, "/")[1:]
// 	}

// 	return routeParts
// }

func ExtractRouteParts(route string) []string {
	return strings.Split(strings.TrimRight(route, "/"), "/")
}

func ServeIndex(ctx *fiber.Ctx) error {
	return writeFile(ctx, "./dist/index.html", fiber.MIMETextHTMLCharsetUTF8)
}

func writeFile(ctx *fiber.Ctx, path string, ctype string) error {
	file, err := os.ReadFile(path)

	if err != nil {
		log.Error(err)

		// Just simply return "not found"
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

func IsSupportedImage(image io.Reader) (bool, error) {
	mtype, err := mimetype.DetectReader(image)
	if err != nil {
		return false, err
	}
	if !slices.Contains(request.IMAGE_FORMATS, mtype.String()) {
		return false, errors.New("uploaded logo must be either a jpeg, jpg or png image")
	}

	return true, nil
}

func GetAbsoluteUrl(ctx *fiber.Ctx) string {
	full_path := ctx.Protocol() + "://" + ctx.Hostname() + ctx.Path()

	if query := string(ctx.Context().URI().QueryString()); query != "" {
		return full_path + "?" + query
	}

	if hash := string(ctx.Context().URI().Hash()); hash != "" {
		return full_path + "#" + hash
	}

	return full_path
}
