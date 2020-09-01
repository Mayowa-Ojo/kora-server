package apiv1

import (
	"github.com/Mayowa-Ojo/kora/config"
	"github.com/Mayowa-Ojo/kora/constants"
	"github.com/Mayowa-Ojo/kora/controllers"
	"github.com/Mayowa-Ojo/kora/repository"
	"github.com/Mayowa-Ojo/kora/services"
	"github.com/gofiber/fiber"
)

// Post - Structure for a post router
type Post struct {
	router fiber.Router
}

// NewPostRouter - Registers all post routes and their respective http handler
func NewPostRouter(br fiber.Router, conn *config.DBConn) {
	router := br.Group("/posts")
	postRepo := repository.NewPostRepository(conn)
	userRepo := repository.NewUserRepository(conn)
	service := services.NewPostService(postRepo, userRepo)
	controller := controllers.NewPostController(service)

	router.Get(constants.GetAllResources, controller.GetAll)
}
