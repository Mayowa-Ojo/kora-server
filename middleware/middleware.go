package middleware

import (
	"github.com/gofiber/fiber"
)

// InitMiddlewares -
func InitMiddlewares(app *fiber.App) {
	app.Use(NotFoundError)
}
