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
	mg "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// CommentService - service layer for comments
type CommentService struct {
	commentRepo domain.CommentRepository
	userRepo    domain.UserRepository
	postRepo    domain.PostRepository
}

// NewCommentService -
func NewCommentService(c domain.CommentRepository, u domain.UserRepository, p domain.PostRepository) domain.CommentService {
	return &CommentService{
		c,
		u,
		p,
	}
}

// GetAll -
func (c *CommentService) GetAll(ctx *fiber.Ctx) ([]entity.Comment, error) {
	postID := ctx.Query("postId")
	slug := ctx.Query("slug") // option to fetch comments through post slug

	if postID != "" {
		postObjectID, err := primitive.ObjectIDFromHex(postID)

		filter := bson.D{{Key: "response_to", Value: postObjectID}}
		opts := options.Find()

		comments, err := c.commentRepo.GetMany(ctx, filter, opts)
		if err != nil {
			if err == mg.ErrNoDocuments {
				return nil, constants.ErrNotFound
			}

			return nil, constants.ErrInternalServer
		}

		return comments, nil
	}

	if slug != "" {
		filter := bson.D{
			{Key: "slug", Value: slug},
			{Key: "post_type", Value: "answer"},
		}
		opts := options.FindOne()

		post, err := c.postRepo.GetOne(ctx, filter, opts)
		if err != nil {
			if err == mg.ErrNoDocuments {
				return nil, constants.ErrNotFound
			}

			return nil, constants.ErrInternalServer
		}

		filter = bson.D{{Key: "response_to", Value: post.ID}}
		findOpts := options.Find()

		comments, err := c.commentRepo.GetMany(ctx, filter, findOpts)
		if err != nil {
			if err == mg.ErrNoDocuments {
				return nil, constants.ErrNotFound
			}

			return nil, constants.ErrInternalServer
		}

		return comments, nil
	}

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

// GetCommentsForPost - get comments for a post that matches given query
func (c *CommentService) GetCommentsForPost(ctx *fiber.Ctx) ([]entity.Comment, error) {
	postID := ctx.Params("id")
	postObjectID, err := primitive.ObjectIDFromHex(postID)

	filter := bson.D{{Key: "response_to", Value: postObjectID}}
	opts := options.Find()

	comments, err := c.commentRepo.GetMany(ctx, filter, opts)
	if err != nil {
		if err == mg.ErrNoDocuments {
			return nil, constants.ErrNotFound
		}

		return nil, constants.ErrInternalServer
	}

	return comments, nil
}
