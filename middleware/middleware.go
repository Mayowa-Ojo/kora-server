package middleware

import (
	"os"

	"github.com/gofiber/cors"
	"github.com/gofiber/fiber"
	"github.com/gofiber/fiber/middleware"
	"github.com/gofiber/helmet"
)

// InitMiddlewares -
func InitMiddlewares(app *fiber.App) {
	app.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:8080"},
		AllowCredentials: true,
	}))

	app.Use(middleware.Logger(middleware.LoggerConfig{
		Format:     "${time} ${method} ${path}",
		TimeFormat: "15:04:05",
		TimeZone:   "Local",
		Output:     os.Stdout,
	}))

	app.Use(middleware.Recover())

	app.Use(helmet.New())
}
