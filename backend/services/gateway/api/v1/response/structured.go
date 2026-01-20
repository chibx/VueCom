package response

import "github.com/gofiber/fiber/v2"

type structuredResponse struct {
	Code    int    `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
}

func NewResponse(ctx *fiber.Ctx, code int, message string, data ...any) error {
	var resp = structuredResponse{
		Code:    code,
		Message: message,
	}

	if len(data) > 0 {
		resp.Data = data[0]
	}

	if code != 0 {
		ctx.Status(code)
	}

	return ctx.JSON(resp)
}

func FromFiberError(ctx *fiber.Ctx, err *fiber.Error, data ...any) error {
	var resp = structuredResponse{
		Code:    err.Code,
		Message: err.Message,
	}

	if len(data) > 0 {
		resp.Data = data[0]
	}

	return ctx.Status(err.Code).JSON(resp)
}
