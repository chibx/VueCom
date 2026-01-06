package middlewares

import (
	"fmt"
	"strconv"
	"vuecom/gateway/internal/constants"
	"vuecom/gateway/internal/types"

	"github.com/go-redis/redis_rate/v10"
	"github.com/gofiber/fiber/v2"
)

func GlobalRateLimit(api *types.Api) fiber.Handler {
	return func(c *fiber.Ctx) error {
		limiter := api.Deps.Limiter
		res, err := limiter.AllowN(c.UserContext(), constants.GlobalLimitKey, constants.GlobalLimit, 1)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Rate limiter error"})
		}
		if res.Allowed == 0 {
			retryAfter := max(int(res.RetryAfter.Seconds()), 1)
			c.Set("Retry-After", strconv.Itoa(retryAfter))
			return c.Status(429).JSON(fiber.Map{"error": "Too many requests (global limit exceeded)"})
		}

		// Optional headers
		c.Set("X-RateLimit-Global-Remaining", strconv.Itoa(res.Remaining))
		return c.Next()
	}
}

func SubRateLimit(api *types.Api) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authType, _ := c.Locals("auth_type").(string)
		rlKey, _ := c.Locals("rl_key").(string)

		// Work on a better way to do this
		var limit redis_rate.Limit
		switch authType {
		case "customer":
			limit = constants.CustomerLimit
		case "admin":
			limit = constants.BackendLimit
		default:
			return c.Next()
		}

		res, err := api.Deps.Limiter.AllowN(c.Context(), rlKey, limit, 1)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Rate limiter error"})
		}
		if res.Allowed == 0 {
			retryAfter := max(int(res.RetryAfter.Seconds()), 1)
			c.Set("Retry-After", strconv.Itoa(retryAfter))
			return c.Status(429).JSON(fiber.Map{
				"error": fmt.Sprintf("Too many requests (%s limit exceeded)", authType),
			})
		}

		// Optional per-client headers
		c.Set("X-RateLimit-Remaining", strconv.Itoa(res.Remaining))
		c.Set("X-RateLimit-Reset", strconv.Itoa(int(res.RetryAfter.Seconds())))
		return c.Next()
	}
}
