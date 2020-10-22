package apiv1

import (
	"github.com/Mayowa-Ojo/kora/config"
	"github.com/Mayowa-Ojo/kora/controllers"
	"github.com/Mayowa-Ojo/kora/repository"
	"github.com/Mayowa-Ojo/kora/services"
	"github.com/gofiber/fiber"
)

// AuthRouter -
type AuthRouter struct {
	router fiber.Router
}

// NewAuthRouter -
func NewAuthRouter(r fiber.Router, conn *config.DBConn) {
	router := r.Group("/auth")
	userRepo := repository.NewUserRepository(conn)
	authService := services.NewAuthService(userRepo)
	authController := controllers.NewAuthController(authService)

	router.Post("/signup", authController.Signup)
	router.Post("/login", authController.Login)
}
