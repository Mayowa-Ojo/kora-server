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

// CommentRepository - data access layer/repository for comments
type CommentRepository struct {
	DB *mongo.Database
}

// NewCommentRepository -
func NewCommentRepository(conn *config.DBConn) domain.CommentRepository {
	return &CommentRepository{conn.DB}
}

// GetAll -
func (c *CommentRepository) GetAll(ctx *fiber.Ctx) ([]entity.Comment, error) {
	col := c.DB.Collection("comments")
	comments := make([]entity.Comment, 0)
	filter := bson.D{{}}

	cur, err := col.Find(ctx.Fasthttp, filter)
	if err != nil {
		return nil, err
	}

	if err := cur.All(ctx.Fasthttp, &comments); err != nil {
		return nil, err
	}

	return comments, nil

}

// GetOne -
func (c *CommentRepository) GetOne(ctx *fiber.Ctx, filter types.Any) (*entity.Comment, error) {
	col := c.DB.Collection("comments")
	comment := new(entity.Comment)

	result := col.FindOne(ctx.Fasthttp, filter)
	if err := result.Decode(&comment); err != nil {
		return nil, err
	}

	return comment, nil
}

// Create -
func (c *CommentRepository) Create(ctx *fiber.Ctx, comment *entity.Comment) (*mongo.InsertOneResult, error) {
	col := c.DB.Collection("comments")

	result, err := col.InsertOne(ctx.Fasthttp, comment)

	if err != nil {
		return nil, err
	}

	return result, nil
}

// GetMany -
func (c *CommentRepository) GetMany(ctx *fiber.Ctx, filter types.Any, opts *options.FindOptions) ([]entity.Comment, error) {
	col := c.DB.Collection("comment")
	comments := make([]entity.Comment, 0)

	cur, err := col.Find(ctx.Fasthttp, filter, opts)
	if err != nil {
		return nil, err
	}

	if err := cur.All(ctx.Fasthttp, &comments); err != nil {
		return nil, err
	}

	return comments, nil
}

// UpdateOne -
func (c *CommentRepository) UpdateOne(ctx *fiber.Ctx, filter, update types.Any) (*mongo.UpdateResult, error) {
	col := c.DB.Collection("comments")

	result, err := col.UpdateOne(ctx.Fasthttp, filter, update)

	if err != nil {
		return nil, err
	}

	return result, nil
}

// DeleteOne -
func (c *CommentRepository) DeleteOne(ctx *fiber.Ctx, filter types.Any) (*mongo.DeleteResult, error) {
	col := c.DB.Collection("comments")

	result, err := col.DeleteOne(ctx.Fasthttp, filter)

	if err != nil {
		return nil, err
	}

	return result, nil
}
