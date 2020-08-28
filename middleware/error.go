package middleware

import "github.com/gofiber/fiber"

// NotFoundError -
func NotFoundError(c *fiber.Ctx) {
	// send 404 error
	err := fiber.NewError(404, "Resource not found.")

	c.Next(err)
}
