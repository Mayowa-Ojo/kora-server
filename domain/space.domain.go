package domain

import (
	"github.com/Mayowa-Ojo/kora/entity"
	"github.com/Mayowa-Ojo/kora/types"
	"github.com/gofiber/fiber"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// SpaceService -
type SpaceService interface {
	GetAll(ctx *fiber.Ctx) ([]entity.Space, error)
	GetOne(ctx *fiber.Ctx) (*entity.Space, error)
	GetBySlug(ctx *fiber.Ctx) (*entity.Space, error)
	Create(ctx *fiber.Ctx) (*entity.Space, error)
	GetPostsForSpace(ctx *fiber.Ctx) ([]entity.Post, error)
	GetMembersForSpace(ctx *fiber.Ctx) (types.GenericMap, error)
	UpdateProfileByAdmin(ctx *fiber.Ctx) (*entity.Space, error)
	DeleteSpaceByAdmin(ctx *fiber.Ctx) error
	FollowSpace(ctx *fiber.Ctx) error
	UnfollowSpace(ctx *fiber.Ctx) error
	SetPinnedPost(ctx *fiber.Ctx) error
	UnsetPinnedPost(ctx *fiber.Ctx) error
}

// SpaceRepository -
type SpaceRepository interface {
	GetAll(ctx *fiber.Ctx) ([]entity.Space, error)
	GetOne(ctx *fiber.Ctx, filter types.Any, opts *options.FindOneOptions) (*entity.Space, error)
	Create(ctx *fiber.Ctx, topic *entity.Space) (*mongo.InsertOneResult, error)
	GetMany(ctx *fiber.Ctx, filter types.Any, opts *options.FindOptions) ([]entity.Space, error)
	UpdateOne(ctx *fiber.Ctx, filter, update types.Any) (*mongo.UpdateResult, error)
	DeleteOne(ctx *fiber.Ctx, filter types.Any) (*mongo.DeleteResult, error)
}
