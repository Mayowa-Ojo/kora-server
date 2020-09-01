package domain

import (
	"github.com/Mayowa-Ojo/kora/entity"
	"github.com/gofiber/fiber"
)

// PostService -
type PostService interface {
	GetAll(ctx *fiber.Ctx) ([]entity.Post, error)
}

// PostRepository -
type PostRepository interface {
	GetAll(ctx *fiber.Ctx) ([]entity.Post, error)
}
