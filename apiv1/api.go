package apiv1

import (
	"github.com/Mayowa-Ojo/kora-server/config"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/gofiber/fiber"
)

// InitRoutes - Setup API all routes
func InitRoutes(app *fiber.App, conn *config.DBConn, sess *session.Session) {
	baseRouter := app.Group("/api/v1")

	NewPostRouter(baseRouter, conn)
	NewAuthRouter(baseRouter, conn)
	NewTopicRouter(baseRouter, conn)
	NewCommentRouter(baseRouter, conn)
	NewUserRouter(baseRouter, conn)
	NewSpaceRouter(baseRouter, conn)
	NewImageRouter(baseRouter, conn, sess)
}
