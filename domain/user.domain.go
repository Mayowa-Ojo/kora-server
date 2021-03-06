package domain

import (
	"github.com/Mayowa-Ojo/kora-server/entity"
	"github.com/Mayowa-Ojo/kora-server/types"
	"github.com/gofiber/fiber"
	mg "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// UserService -
type UserService interface {
	GetAll(ctx *fiber.Ctx) ([]entity.User, error)
	GetOne(ctx *fiber.Ctx) (*entity.User, error)
	GetUserProfile(ctx *fiber.Ctx) (*entity.User, error)
	UpdateContentViews(ctx *fiber.Ctx) error
	UpdateProfile(ctx *fiber.Ctx) (*entity.User, error)
	GetFollowersForUser(ctx *fiber.Ctx) ([]entity.User, error)
	GetFollowingForUser(ctx *fiber.Ctx) ([]entity.User, error)
	GetPostsForUser(ctx *fiber.Ctx) ([]entity.Post, error)
	GetSharedPostsForUser(ctx *fiber.Ctx) ([]entity.SharedPost, error)
	GetSpacesForUser(ctx *fiber.Ctx) ([]entity.Space, error)
	GetKnowledgeForUser(ctx *fiber.Ctx) ([]entity.Topic, error)
	UpdateUserKnowledge(ctx *fiber.Ctx) ([]entity.Topic, error)
	FollowUser(ctx *fiber.Ctx) (*entity.User, error)
	UnfollowUser(ctx *fiber.Ctx) (*entity.User, error)
	SetPinnedPost(ctx *fiber.Ctx) error
	UnsetPinnedPost(ctx *fiber.Ctx) error
}

// UserRepository -
type UserRepository interface {
	GetAll(ctx *fiber.Ctx) ([]entity.User, error)
	GetOne(ctx *fiber.Ctx, filter types.Any, opts *options.FindOneOptions) (*entity.User, error)
	GetMany(ctx *fiber.Ctx, filter types.Any, opts *options.FindOptions) ([]entity.User, error)
	Create(ctx *fiber.Ctx, user entity.User) (*mg.InsertOneResult, error)
	UpdateOne(ctx *fiber.Ctx, filter types.Any, update types.Any) (*mg.UpdateResult, error)
}
