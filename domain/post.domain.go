package domain

import (
	"github.com/Mayowa-Ojo/kora-server/entity"
	"github.com/Mayowa-Ojo/kora-server/types"
	"github.com/gofiber/fiber"
	mg "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// PostService -
type PostService interface {
	GetAll(ctx *fiber.Ctx) ([]entity.Post, error)
	GetOne(ctx *fiber.Ctx) (*entity.Post, error)
	GetBySlug(ctx *fiber.Ctx) (*entity.Post, error)
	Create(ctx *fiber.Ctx) (*entity.Post, error)
	CreateSharedPost(ctx *fiber.Ctx) (*entity.SharedPost, error)
	DeleteOne(ctx *fiber.Ctx) error
	GetFeedForUser(ctx *fiber.Ctx) ([]entity.Post, error)
	GetQuestionsForUser(ctx *fiber.Ctx) ([]entity.Post, error)
	GetAnswersForQuestion(ctx *fiber.Ctx) ([]entity.Post, error)
	GetSuggestedQuestions(ctx *fiber.Ctx) ([]entity.Post, error)
	UpvotePostByUser(ctx *fiber.Ctx) (*entity.Post, error)
	DownvotePostByUser(ctx *fiber.Ctx) (*entity.Post, error)
	FollowPost(ctx *fiber.Ctx) error
	UnfollowPost(ctx *fiber.Ctx) error
	AddTopicsToPost(ctx *fiber.Ctx) ([]entity.Topic, error)
	GetTopicsForPost(ctx *fiber.Ctx) ([]entity.Topic, error)
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

// SharedPostRepository -
type SharedPostRepository interface {
	GetAll(ctx *fiber.Ctx) ([]entity.SharedPost, error)
	GetMany(ctx *fiber.Ctx, filter types.Any, opts *options.FindOptions) ([]entity.SharedPost, error)
	GetOne(ctx *fiber.Ctx, filter types.Any, opts *options.FindOneOptions) (*entity.SharedPost, error)
	UpdateOne(ctx *fiber.Ctx, filter types.Any, update types.Any) (*mg.UpdateResult, error)
	DeleteOne(ctx *fiber.Ctx, filter types.Any) (*mg.DeleteResult, error)
	Create(ctx *fiber.Ctx, sharedPost *entity.SharedPost) (*mg.InsertOneResult, error)
}
