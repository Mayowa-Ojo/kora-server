package repository

import (
	"github.com/Mayowa-Ojo/kora/config"
	"github.com/Mayowa-Ojo/kora/domain"
	"github.com/Mayowa-Ojo/kora/entity"
	"github.com/Mayowa-Ojo/kora/types"
	"github.com/gofiber/fiber"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// TopicRepository - acts as the data access layer/repository for topics
type TopicRepository struct {
	DB *mongo.Database
}

// NewTopicRepository -
func NewTopicRepository(conn *config.DBConn) domain.TopicRepository {
	return &TopicRepository{conn.DB}
}

// GetAll -
func (t *TopicRepository) GetAll(ctx *fiber.Ctx) ([]entity.Topic, error) {
	col := t.DB.Collection("topics")
	topics := make([]entity.Topic, 0)
	filter := bson.D{{}}

	cur, err := col.Find(ctx.Fasthttp, filter)
	if err != nil {
		return nil, err
	}

	if err := cur.All(ctx.Fasthttp, &topics); err != nil {
		return nil, err
	}

	return topics, nil

}

// GetOne -
func (t *TopicRepository) GetOne(ctx *fiber.Ctx, filter types.Any) (*entity.Topic, error) {
	col := t.DB.Collection("topics")
	topic := new(entity.Topic)

	result := col.FindOne(ctx.Fasthttp, filter)
	if err := result.Decode(&topic); err != nil {
		return nil, err
	}

	return topic, nil
}

// Create -
func (t *TopicRepository) Create(ctx *fiber.Ctx, topic *entity.Topic) (*mongo.InsertOneResult, error) {
	col := t.DB.Collection("topics")

	result, err := col.InsertOne(ctx.Fasthttp, topic)

	if err != nil {
		return nil, err
	}

	return result, nil
}

// GetMany -
func (t *TopicRepository) GetMany(ctx *fiber.Ctx, filter types.Any, opts *options.FindOptions) ([]entity.Topic, error) {
	col := t.DB.Collection("topics")
	topics := make([]entity.Topic, 0)

	cur, err := col.Find(ctx.Fasthttp, filter, opts)
	if err != nil {
		return nil, err
	}

	if err := cur.All(ctx.Fasthttp, &topics); err != nil {
		return nil, err
	}

	return topics, nil
}

// UpdateOne -
func (t *TopicRepository) UpdateOne(ctx *fiber.Ctx, filter, update types.Any) (*mongo.UpdateResult, error) {
	col := t.DB.Collection("topics")

	result, err := col.UpdateOne(ctx.Fasthttp, filter, update)

	if err != nil {
		return nil, err
	}

	return result, nil
}

// DeleteOne -
func (t *TopicRepository) DeleteOne(ctx *fiber.Ctx, filter types.Any) (*mongo.DeleteResult, error) {
	col := t.DB.Collection("topics")

	result, err := col.DeleteOne(ctx.Fasthttp, filter)

	if err != nil {
		return nil, err
	}

	return result, nil
}
