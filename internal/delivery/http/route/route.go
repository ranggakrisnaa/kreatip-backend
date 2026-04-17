package route

import "github.com/gofiber/fiber/v2"

// TODO: inject controllers and register routes
func Setup(app *fiber.App) {
	v1 := app.Group("/api/v1")

	// auth
	auth := v1.Group("/auth")
	_ = auth

	// creator profile
	me := v1.Group("/me")
	_ = me

	// public
	creators := v1.Group("/creators")
	_ = creators

	// donations (public)
	donations := v1.Group("/donations")
	_ = donations

	// webhooks (public, no JWT)
	webhooks := v1.Group("/webhooks")
	_ = webhooks

	// admin
	admin := v1.Group("/admin")
	_ = admin
}
