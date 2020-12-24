package services

import (
	"github.com/Mayowa-Ojo/kora/constants"
	"github.com/Mayowa-Ojo/kora/domain"
	"github.com/Mayowa-Ojo/kora/entity"
	"github.com/Mayowa-Ojo/kora/utils"
	"github.com/gofiber/fiber"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// TopicService -
type TopicService struct {
	topicRepo domain.TopicRepository
	postRepo  domain.PostRepository
	userRepo  domain.UserRepository
}

// NewTopicService -
func NewTopicService(t domain.TopicRepository, p domain.PostRepository, u domain.UserRepository) domain.TopicService {
	return &TopicService{
		t,
		p,
		u,
	}
}

// GetAll -
func (t *TopicService) GetAll(ctx *fiber.Ctx) ([]entity.Topic, error) {
	result, err := t.topicRepo.GetAll(ctx)

	if err != nil {
		return nil, constants.ErrInternalServer
	}

	return result, nil
}

// GetOne -
func (t *TopicService) GetOne(ctx *fiber.Ctx) (*entity.Topic, error) {
	id := ctx.Params("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	filter := bson.D{{Key: "_id", Value: objectID}}

	topic, err := t.topicRepo.GetOne(ctx, filter)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, constants.ErrNotFound
		}

		return nil, constants.ErrInternalServer
	}

	return topic, nil
}

// Create -
func (t *TopicService) Create(ctx *fiber.Ctx) (*entity.Topic, error) {
	var requestBody struct {
		Name string `json:"name"`
	}

	err := ctx.BodyParser(&requestBody)
	if err != nil {
		return nil, constants.ErrUnprocessableEntity
	}

	instance := &entity.Topic{
		Name: requestBody.Name,
	}

	if err := instance.Validate(); err != nil {
		return nil, constants.ErrUnprocessableEntity
	}

	instance.SetDefaultValues()

	insertResult, err := t.topicRepo.Create(ctx, instance)
	if err != nil {
		return nil, constants.ErrInternalServer
	}

	filter := bson.D{{Key: "_id", Value: insertResult.InsertedID}}
	topic, err := t.topicRepo.GetOne(ctx, filter)
	if err != nil {
		return nil, constants.ErrInternalServer
	}

	return topic, nil
}

// FollowTopic -
func (t *TopicService) FollowTopic(ctx *fiber.Ctx) error {
	topicID := ctx.Params("id")
	topicObjectID, err := primitive.ObjectIDFromHex(topicID)
	if err != nil {
		return constants.ErrUnprocessableEntity
	}

	user, err := utils.GetUserFromAuthHeader(ctx, t.userRepo)
	if err != nil {
		return constants.ErrUnauthorized
	}

	filter := bson.D{{Key: "_id", Value: topicObjectID}}
	update := bson.D{{Key: "$addToSet", Value: bson.D{{Key: "followers", Value: user.ID}}}}

	if _, err := t.topicRepo.UpdateOne(ctx, filter, update); err != nil {
		return constants.ErrInternalServer
	}

	return nil
}

// UnfollowTopic -
func (t *TopicService) UnfollowTopic(ctx *fiber.Ctx) error {
	topicID := ctx.Params("id")
	topicObjectID, err := primitive.ObjectIDFromHex(topicID)
	if err != nil {
		return constants.ErrUnprocessableEntity
	}

	user, err := utils.GetUserFromAuthHeader(ctx, t.userRepo)
	if err != nil {
		return constants.ErrUnauthorized
	}

	filter := bson.D{{Key: "_id", Value: topicObjectID}}
	update := bson.D{{Key: "$pull", Value: bson.D{{Key: "followers", Value: user.ID}}}}

	if _, err := t.topicRepo.UpdateOne(ctx, filter, update); err != nil {
		return constants.ErrInternalServer
	}

	return nil
}
