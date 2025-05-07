package middlewares

import (
	global_types "chatify/types"
	"chatify/utils"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
)

func RateLimiter(max_requests int, duration time.Duration) fiber.Handler {
	return limiter.New(limiter.Config{
		Max:        max_requests,
		Expiration: duration,
		KeyGenerator: func(ctx *fiber.Ctx) string {
			return utils.GetIPAdress(ctx) // Use the request IP address as the key for rate limiting
		},
		LimitReached: func(ctx *fiber.Ctx) error {
			return ctx.Status(http.StatusTooManyRequests).JSON(global_types.IResponseAPI{
				Message:      "Rate limit exceeded",
				ErrorSection: "RateLimiter | rate limiter",
			})
		},
	})
}
