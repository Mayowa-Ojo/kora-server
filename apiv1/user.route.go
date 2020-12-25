package apiv1

import (
	"github.com/Mayowa-Ojo/kora/config"
	"github.com/Mayowa-Ojo/kora/controllers"
	"github.com/Mayowa-Ojo/kora/middleware"
	"github.com/Mayowa-Ojo/kora/repository"
	"github.com/Mayowa-Ojo/kora/services"
	"github.com/gofiber/fiber"
)

// UserRouter - Structure for a user router
type UserRouter struct {
	router fiber.Router
}

// NewUserRouter - registers all user routes and their respective http handler
// br - base router </api/v1>
// conn - database connection
func NewUserRouter(br fiber.Router, conn *config.DBConn) {
	router := br.Group("/users")
	userRepo := repository.NewUserRepository(conn)
	postRepo := repository.NewPostRepository(conn)
	topicRepo := repository.NewTopicRepository(conn)
	userService := services.NewUserService(userRepo, postRepo, topicRepo)
	userController := controllers.NewUserController(userService)

	router.Get("/", userController.GetAll)
	router.Get("/username", middleware.AuthorizeRoute(), userController.GetUserProfile)
	router.Get("/:id", userController.GetOne)
	router.Patch("/:id", middleware.AuthorizeRoute(), userController.UpdateProfile)
	router.Get("/followers", middleware.AuthorizeRoute(), userController.GetFollowersForUser)
	router.Get("/following", middleware.AuthorizeRoute(), userController.GetFollowingForUser)
	router.Get("/posts", middleware.AuthorizeRoute(), userController.GetPostsForUser)
	router.Get("/knowledge", middleware.AuthorizeRoute(), userController.GetKnowledgeForUser)
	router.Patch("/follow", middleware.AuthorizeRoute(), userController.FollowUser)
	router.Patch("/unfollow", middleware.AuthorizeRoute(), userController.UnfollowUser)
	router.Patch("/pin", middleware.AuthorizeRoute(), userController.SetPinnedPost)
	router.Patch("/unpin", middleware.AuthorizeRoute(), userController.UnsetPinnedPost)
}
