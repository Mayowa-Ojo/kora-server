package domain

import (
	"github.com/Mayowa-Ojo/kora/entity"
	"github.com/Mayowa-Ojo/kora/types"
	"github.com/gofiber/fiber"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// PostService -
type PostService interface {
	GetAll(ctx *fiber.Ctx, opts types.GenericMap) ([]types.GenericMap, error)
}

// PostRepository -
type PostRepository interface {
	GetAll(ctx *fiber.Ctx) ([]entity.Post, error)
	GetMany(ctx *fiber.Ctx, filter types.Any, opts *options.FindOptions) ([]entity.Post, error)
}
