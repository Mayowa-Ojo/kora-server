package repository

import (
	"github.com/Mayowa-Ojo/kora/config"
	"github.com/Mayowa-Ojo/kora/domain"
	"github.com/Mayowa-Ojo/kora/entity"
	"github.com/Mayowa-Ojo/kora/types"
	"github.com/gofiber/fiber"
	"go.mongodb.org/mongo-driver/bson"
	mg "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// SharedPost - acts as the data access layer/repository for posts
type SharedPost struct {
	DB *mg.Database
}

// NewSharedPostRepository -
func NewSharedPostRepository(conn *config.DBConn) domain.SharedPostRepository {
	return &SharedPost{conn.DB}
}

// GetAll -
func (p SharedPost) GetAll(ctx *fiber.Ctx) ([]entity.SharedPost, error) {
	posts := make([]entity.SharedPost, 0)
	filter := bson.D{{}}
	col := p.DB.Collection("shared_posts")

	c, err := col.Find(ctx.Fasthttp, filter)
	if err != nil {
		return nil, err
	}

	if err := c.All(ctx.Fasthttp, &posts); err != nil {
		return nil, err
	}

	return posts, nil
}

// GetOne -
func (p SharedPost) GetOne(ctx *fiber.Ctx, filter types.Any, opts *options.FindOneOptions) (*entity.SharedPost, error) {
	c := p.DB.Collection("shared_posts")
	post := new(entity.SharedPost)

	result := c.FindOne(ctx.Fasthttp, filter, opts)

	if err := result.Decode(&post); err != nil {
		return nil, err
	}

	return post, nil
}

// GetMany -
func (p SharedPost) GetMany(ctx *fiber.Ctx, filter types.Any, opts *options.FindOptions) ([]entity.SharedPost, error) {
	c := p.DB.Collection("shared_posts")
	posts := make([]entity.SharedPost, 0)

	cur, err := c.Find(ctx.Fasthttp, filter, opts)
	if err != nil {
		return nil, err
	}

	if err := cur.All(ctx.Fasthttp, &posts); err != nil {
		return nil, err
	}

	return posts, nil
}

// Create -
func (p SharedPost) Create(ctx *fiber.Ctx, post *entity.SharedPost) (*mg.InsertOneResult, error) {
	c := p.DB.Collection("shared_posts")

	result, err := c.InsertOne(ctx.Fasthttp, post)

	if err != nil {
		return nil, err
	}

	return result, nil
}

// UpdateOne -
func (p SharedPost) UpdateOne(ctx *fiber.Ctx, filter types.Any, update types.Any) (*mg.UpdateResult, error) {
	c := p.DB.Collection("shared_posts")

	result, err := c.UpdateOne(ctx.Fasthttp, filter, update)

	if err != nil {
		return nil, err
	}

	return result, nil
}

// DeleteOne -
func (p SharedPost) DeleteOne(ctx *fiber.Ctx, filter types.Any) (*mg.DeleteResult, error) {
	c := p.DB.Collection("shared_posts")

	result, err := c.DeleteOne(ctx.Fasthttp, filter)

	if err != nil {
		return nil, err
	}

	return result, nil
}
