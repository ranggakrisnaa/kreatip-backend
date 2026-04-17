package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
)

const (
	LocalUserID = "userID"
	LocalRole   = "role"
	LocalJTI    = "jti"
)

// JWT validates the Bearer token and sets userID + role in context locals.
// TODO: inject jwt secret + redis blacklist check
func JWT() fiber.Handler {
	return func(c *fiber.Ctx) error {
		header := c.Get("Authorization")
		if header == "" || !strings.HasPrefix(header, "Bearer ") {
			return fiber.ErrUnauthorized
		}
		// TODO: parse & verify JWT, check Redis blacklist
		return c.Next()
	}
}

// RequireRole checks that the authenticated user has the given role.
func RequireRole(role string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userRole, ok := c.Locals(LocalRole).(string)
		if !ok || userRole != role {
			return fiber.ErrForbidden
		}
		return c.Next()
	}
}

// GetUserID extracts the authenticated user ID from context locals.
func GetUserID(c *fiber.Ctx) string {
	id, _ := c.Locals(LocalUserID).(string)
	return id
}
