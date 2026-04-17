package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/requestid"
)

const HeaderRequestID = "X-Request-Id"

// RequestID injects a unique request ID into every request.
// Fiber's built-in requestid middleware is used; this file just re-exports
// the config for convenience.
func RequestID() fiber.Handler {
	return requestid.New(requestid.Config{
		Header: HeaderRequestID,
	})
}
