package apiv1

import (
	"github.com/gofiber/fiber"
)

// InitRoutes -
func InitRoutes(app *fiber.App) {
	baseRoute := app.Group("/api/v1")

	NewPostRouter(baseRoute)
}
