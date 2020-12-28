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
	sharedPostRepo := repository.NewSharedPostRepository(conn)
	userRepo := repository.NewUserRepository(conn)
	topicRepo := repository.NewTopicRepository(conn)
	commentRepo := repository.NewCommentRepository(conn)
	spaceRepo := repository.NewSpaceRepository(conn)
	postService := services.NewPostService(postRepo, userRepo, sharedPostRepo, commentRepo, spaceRepo, topicRepo)
	userService := services.NewUserService(userRepo, postRepo, topicRepo, spaceRepo, sharedPostRepo)
	controller := controllers.NewPostController(postService, userService)

	router.Get("/feed", middleware.AuthorizeRoute(), controller.GetFeedForUser)
	router.Get("/questions", middleware.AuthorizeRoute(), controller.GetQuestionsForUser)
	router.Get("/suggestions", middleware.AuthorizeRoute(), controller.GetSuggestedQuestions)
	router.Get("/slug", middleware.AuthorizeRoute(), controller.GetBySlug)
	router.Patch("/upvote", middleware.AuthorizeRoute(), controller.UpvotePostByUser)
	router.Patch("/downvote", middleware.AuthorizeRoute(), controller.DownvotePostByUser)
	router.Patch("/follow", middleware.AuthorizeRoute(), controller.FollowPost)
	router.Patch("/unfollow", middleware.AuthorizeRoute(), controller.UnfollowPost)
	router.Post("/:id/share", middleware.AuthorizeRoute(), controller.SharePost)
	router.Post("/:id/topics", middleware.AuthorizeRoute(), controller.AddTopicsToPost)
	router.Get("/:id/topics", middleware.AuthorizeRoute(), controller.GetTopicsForPost)
	router.Get("/:id/answers", middleware.AuthorizeRoute(), controller.GetAnswersForQuestion)
	router.Get("/:id", controller.GetOne)
	router.Get("/", controller.GetAll)
	router.Post("/", middleware.AuthorizeRoute(), controller.Create)
	router.Delete("/:id", controller.DeleteOne)
}
