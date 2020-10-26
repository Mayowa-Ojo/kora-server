package apiv1

import (
	"github.com/Mayowa-Ojo/kora/config"
	"github.com/Mayowa-Ojo/kora/controllers"
	"github.com/Mayowa-Ojo/kora/middleware"
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
	topicRepo := repository.NewTopicRepository(conn)
	postService := services.NewPostService(postRepo, userRepo)
	userService := services.NewUserService(userRepo, postRepo, topicRepo)
	controller := controllers.NewPostController(postService, userService)

	router.Get("/feed", middleware.AuthorizeRoute(), controller.GetFeedForUser)
	router.Patch("/upvote", middleware.AuthorizeRoute(), controller.UpvotePostByUser)
	router.Patch("/downvote", middleware.AuthorizeRoute(), controller.DownvotePostByUser)
	router.Patch("/follow", middleware.AuthorizeRoute(), controller.FollowPost)
	router.Patch("/unfollow", middleware.AuthorizeRoute(), controller.UnfollowPost)
	router.Get("/:id", controller.GetOne)
	router.Get("/", controller.GetAll)
	router.Post("/", middleware.AuthorizeRoute(), controller.Create)
	router.Delete("/:id", controller.DeleteOne)
}
