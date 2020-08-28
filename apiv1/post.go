package apiv1

import (
	"github.com/gofiber/fiber"
)

// Post -
type Post struct {
	router fiber.Router
}

// NewPostRouter -
func NewPostRouter(br fiber.Router) Post {
	router := br.Group("/posts")

	router.Get("/", func(ctx *fiber.Ctx) {
		ctx.JSON(map[string]interface{}{
			"success": true,
			"message": "Welcome to Kora.io",
			"code":    200,
		})
	})

	return Post{
		router,
	}
}
