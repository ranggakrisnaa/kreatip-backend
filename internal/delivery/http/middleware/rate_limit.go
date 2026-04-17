package middleware

import "github.com/gofiber/fiber/v2"

// RateLimit returns a sliding-window rate limit middleware backed by Redis.
// max: max requests per window. windowSec: window size in seconds.
// TODO: implement using go-redis sliding window
func RateLimit(max int, windowSec int) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// TODO: Redis ZADD/ZREMRANGEBYSCORE sliding window per IP
		return c.Next()
	}
}
