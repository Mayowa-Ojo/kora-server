package apiv1

import (
	"github.com/Mayowa-Ojo/kora/config"
	"github.com/Mayowa-Ojo/kora/controllers"
	"github.com/Mayowa-Ojo/kora/repository"
	"github.com/Mayowa-Ojo/kora/services"
	"github.com/gofiber/fiber"
)

// TopicRouter - structure for a topic router
type TopicRouter struct {
	router fiber.Router
}

// NewTopicRouter - Registers all topic routes and their respective http handler
// br - base router </api/v1>
// conn - database connection
func NewTopicRouter(br fiber.Router, conn *config.DBConn) {
	router := br.Group("/topics")
	topicRepo := repository.NewTopicRepository(conn)
	topicService := services.NewTopicService(topicRepo)
	topicController := controllers.NewTopicController(topicService)

	router.Get("/:id", topicController.GetOne)
	router.Get("/", topicController.GetAll)
	router.Post("/", topicController.Create)
}
