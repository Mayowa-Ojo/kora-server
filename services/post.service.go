package services

import (
	"strconv"

	"github.com/Mayowa-Ojo/kora/constants"
	"github.com/Mayowa-Ojo/kora/domain"
	"github.com/Mayowa-Ojo/kora/entity"
	"github.com/Mayowa-Ojo/kora/utils"
	"github.com/gofiber/fiber"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	mg "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// PostService - acts as the business logic layer/service for posts
type PostService struct {
	postRepo domain.PostRepository
	userRepo domain.UserRepository
}

// NewPostService - creates a new post service instance
func NewPostService(p domain.PostRepository, u domain.UserRepository) domain.PostService {
	return &PostService{
		p,
		u,
	}
}

// GetAll - handles business logic to fetch all posts
func (p PostService) GetAll(ctx *fiber.Ctx) ([]entity.Post, error) {
	postType := ctx.Query("postType", "answer")
	limit := ctx.Query("limit", "10")
	limitInt, err := strconv.ParseInt(limit, 10, 64)

	options := options.Find()
	options.SetLimit(limitInt)

	filter := bson.D{{Key: "post_type", Value: postType}}
	result, err := p.postRepo.GetMany(ctx, filter, options)

	if err != nil {
		return nil, constants.ErrInternalServer
	}

	return result, nil
}

// GetOne -
func (p PostService) GetOne(ctx *fiber.Ctx) (*entity.Post, error) {
	id := ctx.Params("id")
	objectID, err := primitive.ObjectIDFromHex(id)

	filter := bson.D{{Key: "_id", Value: objectID}}
	opts := options.FindOne()
	opts.SetProjection(bson.D{{}})
	post, err := p.postRepo.GetOne(ctx, filter, opts)
	if err != nil {
		if err == mg.ErrNoDocuments {
			return nil, constants.ErrNotFound
		}

		return nil, constants.ErrInternalServer
	}

	return post, nil
}

// Create -
func (p PostService) Create(ctx *fiber.Ctx) (*entity.Post, error) {
	var requestBody struct {
		Title       string   `json:"title"`
		Content     string   `json:"content"`
		ContextLink string   `json:"contextLink"`
		PostType    string   `json:"postType"`
		Topics      []string `json:"topics"`
	}

	err := ctx.BodyParser(&requestBody)
	if err != nil {
		return nil, constants.ErrUnprocessableEntity
	}
	instance := &entity.Post{
		PostType: requestBody.PostType,
	}

	switch requestBody.PostType {
	case "question":
		instance.Title = requestBody.Title
		if contextLink := requestBody.ContextLink; contextLink != "" {
			instance.ContextLink = requestBody.ContextLink
		}
	case "answer":
		questionID := ctx.Query("question")
		objectID, err := primitive.ObjectIDFromHex(questionID)
		if err != nil {
			return nil, constants.ErrUnprocessableEntity
		}

		instance.ResponseTo = objectID
		instance.Title = requestBody.Title
		instance.Content = requestBody.Content
	case "post":
		instance.Content = requestBody.Content
	}

	err = instance.Validate()
	if err != nil {
		return nil, constants.ErrUnprocessableEntity
	}

	// fetch author for post
	id, err := utils.GetJwtClaims(ctx, "userId")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, constants.ErrUnauthorized
	}

	filter := bson.D{{Key: "_id", Value: objectID}}
	opts := options.FindOne()
	opts.SetProjection(bson.D{{}})
	user, err := p.userRepo.GetOne(ctx, filter, opts)
	if err != nil {
		return nil, constants.ErrInternalServer
	}

	instance.Author = user

	topicsObjectID := []primitive.ObjectID{}
	for _, v := range requestBody.Topics {
		objectID, err := primitive.ObjectIDFromHex(v)
		if err != nil {
			return nil, constants.ErrUnprocessableEntity
		}
		topicsObjectID = append(topicsObjectID, objectID)
	}

	instance.Topics = topicsObjectID

	insertResult, err := p.postRepo.Create(ctx, instance)
	if err != nil {
		return nil, constants.ErrInternalServer
	}

	filter = bson.D{{Key: "_id", Value: insertResult.InsertedID}}
	opts = options.FindOne()
	opts.SetProjection(bson.D{{}})
	post, err := p.postRepo.GetOne(ctx, filter, opts)
	if err != nil {
		return nil, constants.ErrInternalServer
	}

	return post, nil
}

// GetFeedForUser - fetch posts which satisfies ones of 3 conditions:
//                  <author of the post is followed by the current user>
//                  <the post belongs to a space the current user is subscribed to>
func (p PostService) GetFeedForUser(ctx *fiber.Ctx) ([]entity.Post, error) {
	userID, err := utils.GetJwtClaims(ctx, "userId")
	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, constants.ErrUnauthorized
	}

	filter := bson.D{{Key: "_id", Value: userObjectID}}
	opts := options.FindOne()
	opts.SetProjection(bson.D{{}})
	user, err := p.userRepo.GetOne(ctx, filter, opts)
	if err != nil {
		return nil, constants.ErrInternalServer
	}

	if len(user.Following) < 1 && len(user.Spaces) < 1 {
		return []entity.Post{}, nil
	}

	filter = bson.D{{
		Key: "$or",
		Value: bson.A{
			bson.D{{Key: "author._id", Value: bson.D{{Key: "$in", Value: user.Following}}}},
			bson.D{{Key: "space._id", Value: bson.D{{Key: "$in", Value: user.Spaces}}}},
		},
	}}
	postOpts := options.Find()
	posts, err := p.postRepo.GetMany(ctx, filter, postOpts)

	if err != nil {
		return nil, constants.ErrInternalServer
	}

	return posts, nil
}
