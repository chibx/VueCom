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

// Copied from RequestCtx.StrictFormValue
func StrictFormValue(ctx *fiber.Ctx, key string) string {
	mf, err := ctx.MultipartForm()
	if err == nil && mf.Value != nil {
		vv := mf.Value[key]
		if len(vv) > 0 {
			return vv[0]
		}
	}
	return ""
}

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
		return false, errors.New("Invalid mimetype")
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
