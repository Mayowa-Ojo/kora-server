package services

import (
	"strconv"

	"github.com/Mayowa-Ojo/kora/constants"
	"github.com/Mayowa-Ojo/kora/domain"
	"github.com/Mayowa-Ojo/kora/entity"
	"github.com/Mayowa-Ojo/kora/shorturl"
	"github.com/Mayowa-Ojo/kora/types"
	"github.com/Mayowa-Ojo/kora/utils"
	"github.com/gofiber/fiber"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	mg "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// PostService - acts as the business logic layer/service for posts
type PostService struct {
	postRepo       domain.PostRepository
	userRepo       domain.UserRepository
	sharedPostRepo domain.SharedPostRepository
	commentRepo    domain.CommentRepository
	spaceRepo      domain.SpaceRepository
}

// NewPostService - creates a new post service instance
func NewPostService(
	p domain.PostRepository, u domain.UserRepository, s domain.SharedPostRepository, c domain.CommentRepository, sp domain.SpaceRepository,
) domain.PostService {
	return &PostService{
		p,
		u,
		s,
		c,
		sp,
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

// GetBySlug -
func (p PostService) GetBySlug(ctx *fiber.Ctx) (*entity.Post, error) {
	slug := ctx.Query("slug")
	username := ctx.Query("username")

	filter := bson.D{
		{Key: "slug", Value: slug},
		{Key: "post_type", Value: "question"},
	}
	opts := options.FindOne()
	opts.SetProjection(bson.D{{Key: "author.hash", Value: 0}})

	post, err := p.postRepo.GetOne(ctx, filter, opts)
	if err != nil {
		if err == mg.ErrNoDocuments {
			return nil, constants.ErrNotFound
		}

		return nil, constants.ErrInternalServer
	}

	if username != "" {
		filter = bson.D{
			{Key: "response_to", Value: post.ID},
			{Key: "author.username", Value: username},
		}
		opts = options.FindOne()
		opts.SetProjection(bson.D{{Key: "author.hash", Value: 0}})

		post, err = p.postRepo.GetOne(ctx, filter, opts)
		if err != nil {
			println(err)
			if err == mg.ErrNoDocuments {
				return nil, constants.ErrNotFound
			}

			return nil, constants.ErrInternalServer
		}
	}

	return post, nil
}

// Create -
func (p PostService) Create(ctx *fiber.Ctx) (*entity.Post, error) {
	var requestBody struct {
		Title            string `json:"title"`
		Content          string `json:"content"`
		ContentTruncated string `json:"contentTruncated"`
		ContextLink      string `json:"contextLink"`
		PostType         string `json:"postType"`
	}

	spaceID := ctx.Query("spaceId")

	err := ctx.BodyParser(&requestBody)
	if err != nil {
		return nil, constants.ErrUnprocessableEntity
	}

	slug := utils.GenerateSlug(requestBody.Title)

	instance := &entity.Post{
		PostType: requestBody.PostType,
		Slug:     slug,
	}

	switch requestBody.PostType {
	case "question":
		instance.Title = requestBody.Title
		if contextLink := requestBody.ContextLink; contextLink != "" {
			instance.ContextLink = requestBody.ContextLink
		}
	case "answer":
		questionID := ctx.Query("questionId")
		questionObjectID, err := primitive.ObjectIDFromHex(questionID)
		if err != nil {
			return nil, constants.ErrUnprocessableEntity
		}

		instance.ResponseTo = questionObjectID
		instance.Title = requestBody.Title
		instance.Content = requestBody.Content
		instance.ContentTruncated = requestBody.ContentTruncated
	case "post":
		instance.Content = requestBody.Content
		instance.ContentTruncated = requestBody.ContentTruncated
	}

	instance.SetDefaultValues()

	err = instance.Validate()
	if err != nil {
		return nil, constants.ErrUnprocessableEntity
	}

	// fetch author for post
	user, err := utils.GetUserFromAuthHeader(ctx, p.userRepo)
	if err != nil {
		return nil, err
	}

	url, err := shorturl.CreateURL(ctx, types.GenericMap{
		"slug":     slug,
		"postType": requestBody.PostType,
		"username": user.Username,
	})

	if err != nil {
		return nil, constants.ErrInternalServer
	}

	instance.Author = user
	instance.ShareLink = url.ShortURL

	insertResult, err := p.postRepo.Create(ctx, instance)
	if err != nil {
		return nil, constants.ErrInternalServer
	}

	filter := bson.D{{Key: "_id", Value: insertResult.InsertedID}}
	opts := options.FindOne()
	opts.SetProjection(bson.D{{}})
	post, err := p.postRepo.GetOne(ctx, filter, opts)
	if err != nil {
		return nil, constants.ErrInternalServer
	}

	if spaceID != "" {
		spaceObjectID, err := primitive.ObjectIDFromHex(spaceID)
		if err != nil {
			return nil, constants.ErrInternalServer
		}

		filter := bson.D{{Key: "_id", Value: spaceObjectID}}
		update := bson.D{
			{Key: "$push", Value: bson.D{{Key: "posts", Value: post.ID}}},
		}

		_, err = p.spaceRepo.UpdateOne(ctx, filter, update)
		if err != nil {
			return nil, constants.ErrInternalServer
		}
	}

	return post, nil
}

// CreateSharedPost - create a shared post
func (p PostService) CreateSharedPost(ctx *fiber.Ctx) (*entity.SharedPost, error) {
	postID := ctx.Params("id")
	postObjectID, err := primitive.ObjectIDFromHex(postID)
	if err != nil {
		return nil, constants.ErrUnprocessableEntity
	}

	var requestBody struct {
		ShareComment string `json:"shareComment"`
	}

	err = ctx.BodyParser(&requestBody)
	if err != nil {
		return nil, constants.ErrUnprocessableEntity
	}

	filter := bson.D{{Key: "_id", Value: postObjectID}}
	opts := options.FindOne()
	post, err := p.postRepo.GetOne(ctx, filter, opts)
	if err != nil {
		if err == mg.ErrNoDocuments || err == mg.ErrNilDocument {
			return nil, constants.ErrNotFound
		}
		return nil, constants.ErrInternalServer
	}

	// fetch author for post
	user, err := utils.GetUserFromAuthHeader(ctx, p.userRepo)
	if err != nil {
		return nil, err
	}

	instance := &entity.SharedPost{
		Comment: requestBody.ShareComment,
		Author:  user,
		Post:    post,
	}

	insertResult, err := p.sharedPostRepo.Create(ctx, instance)
	if err != nil {
		return nil, constants.ErrInternalServer
	}

	filter = bson.D{{Key: "_id", Value: insertResult.InsertedID}}
	opts = options.FindOne()
	sharedPost, err := p.sharedPostRepo.GetOne(ctx, filter, opts)
	if err != nil {
		return nil, constants.ErrInternalServer
	}

	return sharedPost, nil
}

// DeleteOne - delete a single post
func (p PostService) DeleteOne(ctx *fiber.Ctx) error {
	postID := ctx.Params("postId")
	postObjectID, err := primitive.ObjectIDFromHex(postID)
	if err != nil {
		return constants.ErrUnprocessableEntity
	}

	filter := bson.D{{Key: "_id", Value: postObjectID}}
	_, err = p.postRepo.DeleteOne(ctx, filter)
	if err != nil {
		return constants.ErrInternalServer
	}

	return nil
}

// GetFeedForUser - fetch posts which satisfies ones of 3 conditions:
//                  <author of the post is followed by the current user>
//                  <the post belongs to a space the current user is subscribed to>
func (p PostService) GetFeedForUser(ctx *fiber.Ctx) ([]entity.Post, error) {
	user, err := utils.GetUserFromAuthHeader(ctx, p.userRepo)
	if err != nil {
		return nil, err
	}

	if len(user.Following) < 1 && len(user.Spaces) < 1 {
		return []entity.Post{}, nil
	}

	filter := bson.D{{
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

// GetQuestionsForUser - fetch questions for current user:
func (p PostService) GetQuestionsForUser(ctx *fiber.Ctx) ([]entity.Post, error) {
	filter := bson.D{{
		Key:   "post_type",
		Value: "question",
	}}
	postOpts := options.Find()
	posts, err := p.postRepo.GetMany(ctx, filter, postOpts)

	if err != nil {
		return nil, constants.ErrInternalServer
	}

	return posts, nil
}

// UpvotePostByUser - upvote a post and add the current user to upvotedBy list
func (p PostService) UpvotePostByUser(ctx *fiber.Ctx) error {
	postID := ctx.Params("postId")
	postObjectID, err := primitive.ObjectIDFromHex(postID)
	if err != nil {
		return constants.ErrUnprocessableEntity
	}

	user, err := utils.GetUserFromAuthHeader(ctx, p.userRepo)
	if err != nil {
		return err
	}

	filter := bson.D{{Key: "_id", Value: postObjectID}}
	update := bson.D{
		{Key: "$inc", Value: bson.D{{Key: "upvotes", Value: 1}}},
		{Key: "$push", Value: bson.D{{Key: "upvoted_by", Value: user.ID}}},
	}

	_, err = p.postRepo.UpdateOne(ctx, filter, update)
	if err != nil {
		return constants.ErrInternalServer
	}

	return nil
}

// DownvotePostByUser - downvote a post and add the current user to downvotedBy list
func (p PostService) DownvotePostByUser(ctx *fiber.Ctx) error {
	postID := ctx.Params("postId")
	postObjectID, err := primitive.ObjectIDFromHex(postID)
	if err != nil {
		return constants.ErrUnprocessableEntity
	}

	user, err := utils.GetUserFromAuthHeader(ctx, p.userRepo)
	if err != nil {
		return err
	}

	filter := bson.D{{Key: "_id", Value: postObjectID}}
	update := bson.D{
		{Key: "$inc", Value: bson.D{{Key: "downvotes", Value: 1}}},
		{Key: "$push", Value: bson.D{{Key: "downvoted_by", Value: user.ID}}},
	}

	_, err = p.postRepo.UpdateOne(ctx, filter, update)
	if err != nil {
		return constants.ErrInternalServer
	}

	return nil
}

// FollowPost - follow a post and add user to followers list
func (p PostService) FollowPost(ctx *fiber.Ctx) error {
	postID := ctx.Params("postId")
	postObjectID, err := primitive.ObjectIDFromHex(postID)
	if err != nil {
		return constants.ErrUnprocessableEntity
	}

	user, err := utils.GetUserFromAuthHeader(ctx, p.userRepo)
	if err != nil {
		return err
	}

	filter := bson.D{{Key: "_id", Value: postObjectID}}
	update := bson.D{{
		Key: "$push", Value: bson.D{{Key: "followers", Value: user.ID}},
	}}

	_, err = p.postRepo.UpdateOne(ctx, filter, update)
	if err != nil {
		return constants.ErrInternalServer
	}

	return nil
}

// UnfollowPost - unfollow a post and remove user from followers list
func (p PostService) UnfollowPost(ctx *fiber.Ctx) error {
	postID := ctx.Params("postId")
	postObjectID, err := primitive.ObjectIDFromHex(postID)
	if err != nil {
		return constants.ErrUnprocessableEntity
	}

	user, err := utils.GetUserFromAuthHeader(ctx, p.userRepo)
	if err != nil {
		return err
	}

	filter := bson.D{{Key: "_id", Value: postObjectID}}
	update := bson.D{{
		Key: "$pull", Value: bson.D{{Key: "followers", Value: user.ID}},
	}}

	_, err = p.postRepo.UpdateOne(ctx, filter, update)
	if err != nil {
		return constants.ErrInternalServer
	}

	return nil
}

// AddTopicToPost - add a new topic to list of topics
func (p PostService) AddTopicToPost(ctx *fiber.Ctx) (*entity.Post, error) {
	var requestBody struct {
		Topics []string `json:"topics"`
	}

	if err := ctx.BodyParser(&requestBody); err != nil {
		return nil, constants.ErrUnprocessableEntity
	}

	topicsObjectID := []primitive.ObjectID{}
	for _, v := range requestBody.Topics {
		objectID, err := primitive.ObjectIDFromHex(v)
		if err != nil {
			return nil, constants.ErrUnprocessableEntity
		}
		topicsObjectID = append(topicsObjectID, objectID)
	}

	postID := ctx.Params("postId")
	postObjectID, err := primitive.ObjectIDFromHex(postID)
	if err != nil {
		return nil, constants.ErrUnprocessableEntity
	}

	filter := bson.D{{Key: "_id", Value: postObjectID}}
	update := bson.D{{
		Key: "$push",
		Value: bson.D{{
			Key:   "topics",
			Value: bson.D{{Key: "$each", Value: topicsObjectID}},
		}},
	}}

	result, err := p.postRepo.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, constants.ErrInternalServer
	}

	filter = bson.D{{Key: "_id", Value: result.UpsertedID}}
	opts := options.FindOne()
	post, err := p.postRepo.GetOne(ctx, filter, opts)
	if err != nil {
		return nil, constants.ErrInternalServer
	}

	return post, nil
}
