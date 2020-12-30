package domain

import (
	"github.com/Mayowa-Ojo/kora/entity"
	"github.com/Mayowa-Ojo/kora/types"
	"github.com/gofiber/fiber"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// CommentService -
type CommentService interface {
	GetAll(ctx *fiber.Ctx) ([]entity.Comment, error)
	GetOne(ctx *fiber.Ctx) (*entity.Comment, error)
	Create(ctx *fiber.Ctx) (*entity.Comment, error)
	CreateCommentReply(ctx *fiber.Ctx) (*entity.Comment, error)
	GetCommentsForPost(ctx *fiber.Ctx) ([]entity.Comment, error)
	AppendCommentsToPost(ctx *fiber.Ctx, post *entity.Post) (*entity.Post, error)
	AppendCommentsToPosts(ctx *fiber.Ctx, posts []entity.Post) ([]entity.Post, error)
}

// CommentRepository -
type CommentRepository interface {
	GetAll(ctx *fiber.Ctx) ([]entity.Comment, error)
	GetOne(ctx *fiber.Ctx, filter types.Any, opts *options.FindOneOptions) (*entity.Comment, error)
	Create(ctx *fiber.Ctx, comment *entity.Comment) (*mongo.InsertOneResult, error)
	GetMany(ctx *fiber.Ctx, filter types.Any, opts *options.FindOptions) ([]entity.Comment, error)
	UpdateOne(ctx *fiber.Ctx, filter, update types.Any) (*mongo.UpdateResult, error)
	DeleteOne(ctx *fiber.Ctx, filter types.Any) (*mongo.DeleteResult, error)
}
