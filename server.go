package main

import (
	"log"

	"github.com/Mayowa-Ojo/kora/apiv1"
	"github.com/Mayowa-Ojo/kora/config"
	"github.com/Mayowa-Ojo/kora/middleware"
	"github.com/Mayowa-Ojo/kora/utils"
	"github.com/gofiber/fiber"
)

func main() {
	app := fiber.New()

	env := config.NewEnvConfig()

	conn, err := config.InitDB(env)
	if err != nil {
		log.Fatal(err)
	}

	apiv1.InitRoutes(app, conn)
	middleware.InitMiddlewares(app)

	app.Settings.ErrorHandler = utils.ErrorHandler

	app.Listen(env.Port)
}
