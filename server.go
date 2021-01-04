package main

import (
	"log"

	"github.com/Mayowa-Ojo/kora-server/shorturl"

	"github.com/Mayowa-Ojo/kora-server/apiv1"
	"github.com/Mayowa-Ojo/kora-server/config"
	"github.com/Mayowa-Ojo/kora-server/middleware"
	"github.com/Mayowa-Ojo/kora-server/utils"
	"github.com/gofiber/fiber"
)

func main() {
	app := fiber.New()

	env := config.NewEnvConfig()

	conn, err := config.InitDB(env)
	if err != nil {
		log.Fatal(err)
	}

	middleware.InitMiddlewares(app)

	sess, err := config.InitAwsSession(env)
	if err != nil {
		log.Fatal(err)
	}

	apiv1.InitRoutes(app, conn, sess)

	shorturl.InitShortURLService(app, conn)

	app.Use(middleware.NotFoundError)

	app.Settings.ErrorHandler = utils.ErrorHandler

	app.Listen(env.Port)
}
