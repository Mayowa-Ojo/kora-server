package domain

import (
	"github.com/Mayowa-Ojo/kora-server/entity"
	"github.com/Mayowa-Ojo/kora-server/types"
	"github.com/gofiber/fiber"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// TopicService -
type TopicService interface {
	GetAll(ctx *fiber.Ctx) ([]entity.Topic, error)
	GetOne(ctx *fiber.Ctx) (*entity.Topic, error)
	Create(ctx *fiber.Ctx) (*entity.Topic, error)
	FollowTopic(ctx *fiber.Ctx) error
	UnfollowTopic(ctx *fiber.Ctx) error
}

// TopicRepository -
type TopicRepository interface {
	GetAll(ctx *fiber.Ctx) ([]entity.Topic, error)
	GetOne(ctx *fiber.Ctx, filter types.Any) (*entity.Topic, error)
	Create(ctx *fiber.Ctx, topic *entity.Topic) (*mongo.InsertOneResult, error)
	GetMany(ctx *fiber.Ctx, filter types.Any, opts *options.FindOptions) ([]entity.Topic, error)
	UpdateOne(ctx *fiber.Ctx, filter, update types.Any) (*mongo.UpdateResult, error)
	DeleteOne(ctx *fiber.Ctx, filter types.Any) (*mongo.DeleteResult, error)
}
