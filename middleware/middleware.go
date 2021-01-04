package middleware

import (
	"os"

	"github.com/Mayowa-Ojo/kora-server/config"
	"github.com/gofiber/cors"
	"github.com/gofiber/fiber"
	"github.com/gofiber/fiber/middleware"
	"github.com/gofiber/helmet"
)

var (
	env = config.NewEnvConfig()
)

// InitMiddlewares -
func InitMiddlewares(app *fiber.App) {
	app.Use(cors.New(cors.Config{
		AllowOrigins:     []string{env.ClientHostname},
		AllowCredentials: true,
	}))

	app.Use(middleware.Logger(middleware.LoggerConfig{
		Format:     "[${time}] ${status} - ${latency} ${method} ${path}\n",
		TimeFormat: "15:04:05",
		TimeZone:   "Local",
		Output:     os.Stdout,
	}))

	app.Use(middleware.Recover())

	app.Use(helmet.New())
}
