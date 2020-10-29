package apiv1

import (
	"github.com/Mayowa-Ojo/kora/config"
	"github.com/Mayowa-Ojo/kora/controllers"
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

	router.Get("/:id", spaceController.GetOne)
	router.Get("/", spaceController.GetAll)
	router.Post("/", spaceController.Create)
	router.Get("/:id/posts", spaceController.GetPostsForSpace)
	router.Get("/:id/people", spaceController.GetMembersForSpace)
	router.Patch("/:id", spaceController.UpdateProfileByAdmin)
	router.Patch("/:/follow", spaceController.FollowSpace)
	router.Patch("/:id/unfollow", spaceController.UnfollowSpace)
	router.Patch("/:id/pin", spaceController.SetPinnedPost)
	router.Patch("/:id/unpin", spaceController.UnsetPinnedPost)
	router.Delete("/:id", spaceController.DeleteSpaceByAdmin)
}
