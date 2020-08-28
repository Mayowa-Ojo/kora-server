package main

import (
	"github.com/Mayowa-Ojo/kora/apiv1"
	"github.com/Mayowa-Ojo/kora/middleware"
	"github.com/Mayowa-Ojo/kora/utils"
	"github.com/gofiber/fiber"
)

func main() {
	// conn, err := config.Connect()

	// if err != nil {
	// 	log.Fatal(err)
	// }

	app := fiber.New()

	apiv1.InitRoutes(app)
	middleware.InitMiddlewares(app)

	app.Settings.ErrorHandler = utils.ErrorHandler

	// bookRouter := api.Group("/books")
	// userRouter := api.Group("/users")

	// bookRepository := book.NewRepository(conn)
	// userRepository := user.NewRepository(conn)

	// bookService := book.NewService(bookRepository, userRepository)
	// userService := user.NewService(userRepository)

	// book.NewController(bookService, bookRouter)
	// user.NewController(userService, userRouter)

	// app.Use(utils.NotFoundError)

	app.Listen(4000)
}
