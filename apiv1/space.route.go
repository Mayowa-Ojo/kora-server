package apiv1

import (
	"github.com/Mayowa-Ojo/kora/config"
	"github.com/Mayowa-Ojo/kora/controllers"
	"github.com/Mayowa-Ojo/kora/middleware"
	"github.com/Mayowa-Ojo/kora/repository"
	"github.com/Mayowa-Ojo/kora/services"
	"github.com/gofiber/fiber"
)

// SpaceRouter - structure for a topic router
type SpaceRouter struct {
	router fiber.Router
}

// NewSpaceRouter - Registers all topic routes and their respective http handler
// br - base router </api/v1>
// conn - database connection
func NewSpaceRouter(br fiber.Router, conn *config.DBConn) {
	router := br.Group("/spaces")
	spaceRepo := repository.NewSpaceRepository(conn)
	postRepo := repository.NewPostRepository(conn)
	userRepo := repository.NewUserRepository(conn)
	spaceService := services.NewSpaceService(spaceRepo, postRepo, userRepo)
	spaceController := controllers.NewSpaceController(spaceService)

	router.Get("/slug", middleware.AuthorizeRoute(), spaceController.GetBySlug)
	router.Get("/suggestions", middleware.AuthorizeRoute(), spaceController.GetSuggestedSpaces)
	router.Get("/:id", middleware.AuthorizeRoute(), spaceController.GetOne)
	router.Get("/", middleware.AuthorizeRoute(), spaceController.GetAll)
	router.Post("/", middleware.AuthorizeRoute(), spaceController.Create)
	router.Get("/:id/posts", middleware.AuthorizeRoute(), spaceController.GetPostsForSpace)
	router.Get("/:id/people", middleware.AuthorizeRoute(), spaceController.GetMembersForSpace)
	router.Patch("/:id", middleware.AuthorizeRoute(), spaceController.UpdateProfileByAdmin)
	router.Patch("/:/follow", middleware.AuthorizeRoute(), spaceController.FollowSpace)
	router.Patch("/:id/unfollow", middleware.AuthorizeRoute(), spaceController.UnfollowSpace)
	router.Patch("/:id/pin", middleware.AuthorizeRoute(), spaceController.SetPinnedPost)
	router.Patch("/:id/unpin", middleware.AuthorizeRoute(), spaceController.UnsetPinnedPost)
	router.Delete("/:id", middleware.AuthorizeRoute(), spaceController.DeleteSpaceByAdmin)
}
