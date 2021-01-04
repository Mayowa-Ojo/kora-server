package services

import (
	"fmt"

	"github.com/Mayowa-Ojo/kora-server/constants"
	"github.com/Mayowa-Ojo/kora-server/domain"
	"github.com/Mayowa-Ojo/kora-server/entity"
	"github.com/Mayowa-Ojo/kora-server/utils"
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
	opts := options.FindOne()

	comment, err := c.commentRepo.GetOne(ctx, filter, opts)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, constants.ErrNotFound
		}

		return nil, constants.ErrInternalServer
	}

	return comment, nil
}

// Create - [POST] </comments?postId="">
func (c *CommentService) Create(ctx *fiber.Ctx) (*entity.Comment, error) {
	var requestBody struct {
		Content string `json:"content"`
	}

	postID := ctx.Query("postId")
	postObjectID, err := primitive.ObjectIDFromHex(postID)
	if err != nil {
		return nil, constants.ErrUnprocessableEntity
	}

	if err := ctx.BodyParser(&requestBody); err != nil {
		return nil, constants.ErrUnprocessableEntity
	}

	user, err := utils.GetUserFromAuthHeader(ctx, c.userRepo)
	if err != nil {
		return nil, err
	}

	instance := &entity.Comment{
		Content:    requestBody.Content,
		Author:     user,
		ResponseTo: postObjectID,
	}

	instance.SetDefaultValues()

	if err := instance.Validate(); err != nil {
		fmt.Println(err)
		return nil, constants.ErrUnprocessableEntity
	}

	insertResult, err := c.commentRepo.Create(ctx, instance)
	if err != nil {
		return nil, constants.ErrInternalServer
	}

	filter := bson.D{{Key: "_id", Value: insertResult.InsertedID}}
	opts := options.FindOne()

	comment, err := c.commentRepo.GetOne(ctx, filter, opts)
	if err != nil {
		return nil, constants.ErrInternalServer
	}

	return comment, nil
}

// CreateCommentReply - create a reply to a comment that matches given id
// [POST] </comments/reply?commentId="">
func (c *CommentService) CreateCommentReply(ctx *fiber.Ctx) (*entity.Comment, error) {
	var requestBody struct {
		Content string `json:"content"`
	}

	commentID := ctx.Query("commentId")
	commentObjectID, err := primitive.ObjectIDFromHex(commentID)
	if err != nil {
		return nil, constants.ErrUnprocessableEntity
	}

	err = ctx.BodyParser(&requestBody)
	if err != nil {
		return nil, constants.ErrUnprocessableEntity
	}

	user, err := utils.GetUserFromAuthHeader(ctx, c.userRepo)
	if err != nil {
		return nil, err
	}

	instance := &entity.Comment{
		Content:    requestBody.Content,
		Author:     user,
		ResponseTo: commentObjectID,
	}

	if err := instance.Validate(); err != nil {
		return nil, constants.ErrUnprocessableEntity
	}

	insertResult, err := c.commentRepo.Create(ctx, instance)
	if err != nil {
		return nil, constants.ErrInternalServer
	}

	filter := bson.D{{Key: "_id", Value: insertResult.InsertedID}}
	opts := options.FindOne()

	commentReply, err := c.commentRepo.GetOne(ctx, filter, opts)
	if err != nil {
		return nil, constants.ErrInternalServer
	}

	filter = bson.D{{Key: "_id", Value: commentObjectID}}
	update := bson.D{{Key: "$addToSet", Value: bson.D{{Key: "replies", Value: commentReply}}}}

	if _, err := c.commentRepo.UpdateOne(ctx, filter, update); err != nil {
		return nil, constants.ErrInternalServer
	}

	return commentReply, nil
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

// UpvoteCommentByUser - upvote a comment
func (c *CommentService) UpvoteCommentByUser(ctx *fiber.Ctx) (*entity.Comment, error) {
	commentID := ctx.Params("id")
	commentObjectID, err := primitive.ObjectIDFromHex(commentID)
	if err != nil {
		return nil, constants.ErrUnprocessableEntity
	}

	filter := bson.D{{Key: "_id", Value: commentObjectID}}
	update := bson.D{
		{Key: "$inc", Value: bson.D{{Key: "upvotes", Value: 1}}},
	}

	if _, err := c.commentRepo.UpdateOne(ctx, filter, update); err != nil {
		return nil, constants.ErrInternalServer
	}

	filter = bson.D{{Key: "_id", Value: commentObjectID}}
	opts := options.FindOne()
	opts.SetProjection(bson.D{{Key: "author.hash", Value: 0}})

	comment, err := c.commentRepo.GetOne(ctx, filter, opts)
	if err != nil {
		fmt.Println(err)
		return nil, constants.ErrInternalServer
	}

	if comment.Replies == nil {
		filter := bson.D{
			{Key: "_id", Value: comment.ResponseTo},
			{Key: "replies._id", Value: commentObjectID},
		}
		update := bson.D{
			{Key: "$inc", Value: bson.D{{Key: "replies.$.upvotes", Value: 1}}},
		}

		if _, err := c.commentRepo.UpdateOne(ctx, filter, update); err != nil {
			return nil, constants.ErrInternalServer
		}
	}

	return comment, nil
}

// DownvoteCommentByUser - upvote a comment
func (c *CommentService) DownvoteCommentByUser(ctx *fiber.Ctx) (*entity.Comment, error) {
	commentID := ctx.Params("id")
	commentObjectID, err := primitive.ObjectIDFromHex(commentID)
	if err != nil {
		return nil, constants.ErrUnprocessableEntity
	}

	filter := bson.D{{Key: "_id", Value: commentObjectID}}
	update := bson.D{
		{Key: "$inc", Value: bson.D{{Key: "downvotes", Value: 1}}},
	}

	if _, err := c.commentRepo.UpdateOne(ctx, filter, update); err != nil {
		return nil, constants.ErrInternalServer
	}

	filter = bson.D{{Key: "_id", Value: commentObjectID}}
	opts := options.FindOne()

	comment, err := c.commentRepo.GetOne(ctx, filter, opts)
	if err != nil {
		fmt.Println(err)
		return nil, constants.ErrInternalServer
	}

	return comment, nil
}

// AppendCommentsToPost - populate comments field for a given post
func (c *CommentService) AppendCommentsToPost(ctx *fiber.Ctx, post *entity.Post) (*entity.Post, error) {
	filter := bson.D{{Key: "response_to", Value: post.ID}}
	opts := options.Find()

	comments, err := c.commentRepo.GetMany(ctx, filter, opts)
	if err != nil {
		return nil, constants.ErrInternalServer
	}
	fmt.Println(len(comments))
	post.Comments = comments

	return post, nil
}

// AppendCommentsToPosts - populate comments field for a group of post
func (c *CommentService) AppendCommentsToPosts(ctx *fiber.Ctx, posts []entity.Post) ([]entity.Post, error) {
	var filter bson.D
	opts := options.Find()

	for _, v := range posts {
		filter = bson.D{{Key: "response_to", Value: v.ID}}
		comments, err := c.commentRepo.GetMany(ctx, filter, opts)

		if err != nil {
			return nil, constants.ErrInternalServer
		}

		v.Comments = comments
	}

	return posts, nil
}
