package apiv1

import (
	"github.com/Mayowa-Ojo/kora/config"
	"github.com/gofiber/fiber"
)

// InitRoutes - Setup API all routes
func InitRoutes(app *fiber.App, conn *config.DBConn) {
	baseRoute := app.Group("/api/v1")

	NewPostRouter(baseRoute, conn)
}
