package domain

import (
	"github.com/Mayowa-Ojo/kora/entity"
	"github.com/Mayowa-Ojo/kora/types"
	"github.com/gofiber/fiber"
	mg "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// PostService -
type PostService interface {
	GetAll(ctx *fiber.Ctx) ([]entity.Post, error)
	GetOne(ctx *fiber.Ctx) (*entity.Post, error)
	Create(ctx *fiber.Ctx) (*entity.Post, error)
	DeleteOne(ctx *fiber.Ctx) error
	GetFeedForUser(ctx *fiber.Ctx) ([]entity.Post, error)
	UpvotePostByUser(ctx *fiber.Ctx) error
	DownvotePostByUser(ctx *fiber.Ctx) error
	FollowPost(ctx *fiber.Ctx) error
	UnfollowPost(ctx *fiber.Ctx) error
	// GetMany(ctx *fiber.Ctx) (*entity.Post, error)
	// UpdateOne(ctx *fiber.Ctx) (*entity.Post, error)
	// UpdateMany(ctx *fiber.Ctx) (*entity.Post, error)
	// DeleteOne(ctx *fiber.Ctx) (*entity.Post, error)
}

// PostRepository -
type PostRepository interface {
	GetAll(ctx *fiber.Ctx) ([]entity.Post, error)
	GetMany(ctx *fiber.Ctx, filter types.Any, opts *options.FindOptions) ([]entity.Post, error)
	GetOne(ctx *fiber.Ctx, filter types.Any, opts *options.FindOneOptions) (*entity.Post, error)
	UpdateOne(ctx *fiber.Ctx, filter types.Any, update types.Any) (*mg.UpdateResult, error)
	DeleteOne(ctx *fiber.Ctx, filter types.Any) (*mg.DeleteResult, error)
	Create(ctx *fiber.Ctx, post *entity.Post) (*mg.InsertOneResult, error)
}
