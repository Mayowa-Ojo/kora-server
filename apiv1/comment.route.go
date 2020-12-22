package apiv1

import (
	"github.com/Mayowa-Ojo/kora/config"
	"github.com/Mayowa-Ojo/kora/controllers"
	"github.com/Mayowa-Ojo/kora/middleware"
	"github.com/Mayowa-Ojo/kora/repository"
	"github.com/Mayowa-Ojo/kora/services"
	"github.com/gofiber/fiber"
)

// CommentRouter - Structure for a comment router
type CommentRouter struct {
	router fiber.Router
}

// NewCommentRouter - Registers all comment routes and their respective http handler
// br - base router </api/v1>
// conn - database connection
func NewCommentRouter(br fiber.Router, conn *config.DBConn) {
	router := br.Group("/comments")
	commentRepo := repository.NewCommentRepository(conn)
	userRepo := repository.NewUserRepository(conn)
	postRepo := repository.NewPostRepository(conn)
	commentService := services.NewCommentService(commentRepo, userRepo, postRepo)
	commentController := controllers.NewCommentController(commentService)

	router.Get("/:id", commentController.GetOne)
	router.Get("/", commentController.GetAll)
	router.Post("/", middleware.AuthorizeRoute(), commentController.Create) // </comments?postId={}>
}
