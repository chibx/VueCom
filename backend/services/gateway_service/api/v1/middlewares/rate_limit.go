package middlewares

import (
	"fmt"
	"strconv"
	"vuecom/gateway/api/v1/response"
	"vuecom/gateway/internal/constants"
	"vuecom/gateway/internal/types"

	"github.com/go-redis/redis_rate/v10"
	"github.com/gofiber/fiber/v2"
)

func GlobalRateLimit(api *types.Api) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		limiter := api.Deps.Limiter
		res, err := limiter.AllowN(ctx.UserContext(), constants.GlobalLimitKey, constants.GlobalLimit, 1)
		if err != nil {
			return response.NewResponse(ctx, fiber.StatusInternalServerError, "", fiber.Map{"error": "Rate limiter error"})
		}
		if res.Allowed == 0 {
			retryAfter := max(int(res.RetryAfter.Seconds()), 1)
			ctx.Set("Retry-After", strconv.Itoa(retryAfter))
			return response.NewResponse(ctx, fiber.StatusTooManyRequests, "", fiber.Map{"error": "Too many requests (global limit exceeded)"})
		}

		// Optional headers
		ctx.Set("X-RateLimit-Global-Limit", strconv.Itoa(constants.GlobalLimit.Rate))
		ctx.Set("X-RateLimit-Global-Reset", strconv.Itoa(int(res.RetryAfter.Seconds())))
		ctx.Set("X-RateLimit-Global-Remaining", strconv.Itoa(res.Remaining))
		return ctx.Next()
	}
}

func SubRateLimit(api *types.Api) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		authType, _ := ctx.Locals("auth_type").(string)
		rlKey, _ := ctx.Locals("rl_key").(string)

		// Work on a better way to do this
		var limit redis_rate.Limit
		switch authType {
		case "customer":
			limit = constants.CustomerLimit
		case "admin":
			limit = constants.BackendLimit
		default:
			return ctx.Next()
		}

		res, err := api.Deps.Limiter.AllowN(ctx.Context(), rlKey, limit, 1)
		if err != nil {
			return response.NewResponse(ctx, fiber.StatusInternalServerError, "", fiber.Map{"error": "Rate limiter error"})
		}
		if res.Allowed == 0 {
			retryAfter := max(int(res.RetryAfter.Seconds()), 1)
			ctx.Set("Retry-After", strconv.Itoa(retryAfter))
			return response.NewResponse(ctx, fiber.StatusTooManyRequests, "", fiber.Map{
				"error": fmt.Sprintf("Too many requests (%s limit exceeded)", authType),
			})
		}

		// Optional per-client headers
		ctx.Set("X-RateLimit-Limit", strconv.Itoa(limit.Rate))
		ctx.Set("X-RateLimit-Remaining", strconv.Itoa(res.Remaining))
		ctx.Set("X-RateLimit-Reset", strconv.Itoa(int(res.RetryAfter.Seconds())))
		return ctx.Next()
	}
}
