package middlewares

import (
	"strconv"

	"github.com/chibx/vuecom/backend/services/gateway/api/v1/response"
	"github.com/chibx/vuecom/backend/services/gateway/internal/constants"
	"github.com/chibx/vuecom/backend/services/gateway/internal/types"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func GlobalRateLimit(api *types.Api) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		logger := api.Deps.Logger
		res, err := api.Deps.Limiter.Allow(ctx.UserContext(), constants.GlobalLimitKey, constants.GlobalLimit)
		if err != nil {
			logger.Error("failed to allow global rate limit", zap.Error(err))
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

func BackendRateLimit(api *types.Api) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		logger := api.Deps.Logger
		rlKey, _ := ctx.Locals("rl_key").(string)

		// Work on a better way to do this
		var limit = constants.BackendLimit

		res, err := api.Deps.Limiter.Allow(ctx.Context(), rlKey, limit)
		if err != nil {
			logger.Error("failed to allow backend rate limit", zap.Error(err))
			return response.NewResponse(ctx, fiber.StatusInternalServerError, "", fiber.Map{"error": "Rate limiter error"})
		}
		if res.Allowed == 0 {
			retryAfter := max(int(res.RetryAfter.Seconds()), 1)
			ctx.Set("Retry-After", strconv.Itoa(retryAfter))
			return response.NewResponse(ctx, fiber.StatusTooManyRequests, "", fiber.Map{
				"error": "Too many requests (backend limit exceeded)",
			})
		}

		// Optional per-client headers
		ctx.Set("X-RateLimit-Limit", strconv.Itoa(limit.Rate))
		ctx.Set("X-RateLimit-Remaining", strconv.Itoa(res.Remaining))
		ctx.Set("X-RateLimit-Reset", strconv.Itoa(int(res.RetryAfter.Seconds())))
		return ctx.Next()
	}
}

func CustomerRateLimit(api *types.Api) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		logger := api.Deps.Logger
		customerID := ctx.Get(constants.CustomerHeaderKey)
		var limit = constants.CustomerLimit // Could maybe add a fallback for anonymous users but it aint compulsory
		var rlKey string
		if customerID != "" {
			rlKey = constants.CustomerLimitKey + customerID
		} else {
			rlKey = constants.AnonymousLimitKey + ctx.IP()
		}

		// Work on a better way to do this

		res, err := api.Deps.Limiter.Allow(ctx.Context(), rlKey, limit)
		if err != nil {
			logger.Error("failed to allow customer rate limit", zap.Error(err))
			return response.NewResponse(ctx, fiber.StatusInternalServerError, "", fiber.Map{"error": "Rate limiter error"})
		}
		if res.Allowed == 0 {
			retryAfter := max(int(res.RetryAfter.Seconds()), 1)
			ctx.Set("Retry-After", strconv.Itoa(retryAfter))
			return response.NewResponse(ctx, fiber.StatusTooManyRequests, "", fiber.Map{
				"error": "Too many requests (customer limit exceeded)",
			})
		}

		// Optional per-client headers
		ctx.Set("X-RateLimit-Limit", strconv.Itoa(limit.Rate))
		ctx.Set("X-RateLimit-Remaining", strconv.Itoa(res.Remaining))
		ctx.Set("X-RateLimit-Reset", strconv.Itoa(int(res.RetryAfter.Seconds())))
		return ctx.Next()
	}
}
