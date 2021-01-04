package apiv1

import (
	"github.com/Mayowa-Ojo/kora-server/config"
	"github.com/Mayowa-Ojo/kora-server/controllers"
	"github.com/Mayowa-Ojo/kora-server/repository"
	"github.com/Mayowa-Ojo/kora-server/services"
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
