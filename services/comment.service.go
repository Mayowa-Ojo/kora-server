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
	"go.mongodb.org/mongo-driver/mongo/options"
)

// CommentService - service layer for comments
type CommentService struct {
	commentRepo domain.CommentRepository
	userRepo    domain.UserRepository
}

// NewCommentService -
func NewCommentService(c domain.CommentRepository, u domain.UserRepository) domain.CommentService {
	return &CommentService{
		c,
		u,
	}
}

// GetAll -
func (c *CommentService) GetAll(ctx *fiber.Ctx) ([]entity.Comment, error) {
	result, err := c.commentRepo.GetAll(ctx)

	if err != nil {
		return nil, constants.ErrInternalServer
	}

	return result, nil
}

// GetOne -
func (c *CommentService) GetOne(ctx *fiber.Ctx) (*entity.Comment, error) {
	id := ctx.Params("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	filter := bson.D{{Key: "_id", Value: objectID}}

	comment, err := c.commentRepo.GetOne(ctx, filter)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, constants.ErrNotFound
		}

		return nil, constants.ErrInternalServer
	}

	return comment, nil
}

// Create -
func (c *CommentService) Create(ctx *fiber.Ctx) (*entity.Comment, error) {
	var requestBody struct {
		Content string `json:"content"`
	}
	postID := ctx.Query("postId")
	postObjectID, err := primitive.ObjectIDFromHex(postID)
	if err != nil {
		return nil, constants.ErrUnprocessableEntity
	}

	err = ctx.BodyParser(&requestBody)
	if err != nil {
		return nil, constants.ErrUnprocessableEntity
	}

	instance := &entity.Comment{
		Content:    requestBody.Content,
		ResponseTo: postObjectID,
	}

	userID, err := utils.GetJwtClaims(ctx, "userId")
	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, constants.ErrUnauthorized
	}

	filter := bson.D{{Key: "_id", Value: userObjectID}}
	opts := options.FindOne()
	opts.SetProjection(bson.D{{}})
	user, err := c.userRepo.GetOne(ctx, filter, opts)
	if err != nil {
		return nil, constants.ErrInternalServer
	}

	instance.Author = user

	if err := instance.Validate(); err != nil {
		return nil, constants.ErrUnprocessableEntity
	}

	insertResult, err := c.commentRepo.Create(ctx, instance)
	if err != nil {
		return nil, constants.ErrInternalServer
	}

	filter = bson.D{{Key: "_id", Value: insertResult.InsertedID}}
	comment, err := c.commentRepo.GetOne(ctx, filter)
	if err != nil {
		return nil, constants.ErrInternalServer
	}

	return comment, nil
}
