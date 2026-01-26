package response

import "github.com/gofiber/fiber/v2"

type ErrorDetail struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type structuredResponse struct {
	Code    int    `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
	Errors  any    `json:"errors,omitempty"`
	Data    any    `json:"data,omitempty"`
}

func isSlice(v any) bool {
	_, ok := v.([]any) // Or customize for your error types
	return ok
}

func isErrorField(v any) bool {
	_, ok := v.([]ErrorDetail) // Or customize for your error types
	// return ok
	return ok
}

func NewResponse(code int, message string, data ...any) *structuredResponse {
	if code == 0 {
		code = fiber.StatusOK
	}

	var resp = &structuredResponse{
		Code:    code,
		Message: message,
	}

	if len(data) > 0 {
		// If first data is a ErrorDetail, assume it's errors; else, it's Data
		if errs, ok := data[0].([]ErrorDetail); ok {
			resp.Errors = errs
		} else {
			resp.Data = data[0]
		}
	}

	return resp
}

func From(ctx *fiber.Ctx, structured *structuredResponse) error {
	if structured == nil {
		return nil
	}

	if structured.Code != 0 {
		ctx.Status(structured.Code)
	}

	return ctx.JSON(structured)
}

func WriteResponse(ctx *fiber.Ctx, code int, message string, data ...any) error {
	var resp = structuredResponse{
		Code:    code,
		Message: message,
	}

	if len(data) > 0 {
		// If first data is a ErrorDetail, assume it's errors; else, it's Data
		if errs, ok := data[0].([]ErrorDetail); ok {
			resp.Errors = errs
		} else {
			resp.Data = data[0]
		}
	}

	if code != 0 {
		ctx.Status(code)
	} else {
		ctx.Status(fiber.StatusOK)
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
