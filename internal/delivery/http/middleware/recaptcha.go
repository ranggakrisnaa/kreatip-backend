package middleware

import "github.com/gofiber/fiber/v2"

// ReCaptcha validates the g-recaptcha-response token via Google's verify API.
// TODO: inject site secret key + min score threshold
func ReCaptcha() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// TODO: read token from header X-Recaptcha-Token, call Google API
		return c.Next()
	}
}
