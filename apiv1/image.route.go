package apiv1

import (
	"github.com/Mayowa-Ojo/kora-server/config"
	"github.com/Mayowa-Ojo/kora-server/controllers"
	"github.com/Mayowa-Ojo/kora-server/repository"
	"github.com/Mayowa-Ojo/kora-server/services"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/gofiber/fiber"
)

// ImageRouter -
type ImageRouter struct {
	router fiber.Router
}

// NewImageRouter -
func NewImageRouter(r fiber.Router, conn *config.DBConn, sess *session.Session) {
	router := r.Group("/images")
	userRepo := repository.NewUserRepository(conn)
	postRepo := repository.NewPostRepository(conn)
	spaceRepo := repository.NewSpaceRepository(conn)
	imageService := services.NewImageService(postRepo, userRepo, spaceRepo, sess)
	imageController := controllers.NewImageController(imageService)

	router.Post("/", imageController.UploadImage)
}
