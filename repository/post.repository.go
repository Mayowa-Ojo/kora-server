package repository

import (
	"github.com/Mayowa-Ojo/kora/config"
	"github.com/Mayowa-Ojo/kora/domain"
	"github.com/Mayowa-Ojo/kora/entity"
	"github.com/Mayowa-Ojo/kora/types"
	"github.com/gofiber/fiber"
	"go.mongodb.org/mongo-driver/bson"
	mg "go.mongodb.org/mongo-driver/mongo"
)

// Post - acts as the data access layer/repository for posts
type Post struct {
	DB *mg.Database
}

// NewPostRepository -
func NewPostRepository(conn *config.DBConn) domain.PostRepository {
	return &Post{conn.DB}
}

// GetAll -
func (p Post) GetAll(ctx *fiber.Ctx) ([]entity.Post, error) {
	posts := make([]entity.Post, 0)
	filter := bson.D{{}}
	col := p.DB.Collection("posts")

	c, err := col.Find(ctx.Fasthttp, filter)
	if err != nil {
		return nil, err
	}

	if err := c.All(ctx.Fasthttp, &posts); err != nil {
		return nil, err
	}

	return posts, nil
}

// GetByID -
func (p Post) GetByID(ctx *fiber.Ctx, filter types.Any) (*entity.Post, error) {
	c := p.DB.Collection("posts")
	post := new(entity.Post)

	result := c.FindOne(ctx.Fasthttp, filter)

	if err := result.Decode(&post); err != nil {
		return nil, err
	}

	return post, nil
}

// Create -
func (p Post) Create(ctx *fiber.Ctx, post *entity.Post) (*mg.InsertOneResult, error) {
	c := p.DB.Collection("posts")

	result, err := c.InsertOne(ctx.Fasthttp, post)

	if err != nil {
		return nil, err
	}

	return result, nil
}

// UpdateOne -
func (p Post) UpdateOne(ctx *fiber.Ctx, filter types.Any, update types.Any) (*mg.UpdateResult, error) {
	c := p.DB.Collection("posts")

	result, err := c.UpdateOne(ctx.Fasthttp, filter, update)

	if err != nil {
		return nil, err
	}

	return result, nil
}

// DeleteOne -
func (p Post) DeleteOne(ctx *fiber.Ctx, filter types.Any) (*mg.DeleteResult, error) {
	c := p.DB.Collection("posts")

	result, err := c.DeleteOne(ctx.Fasthttp, filter)

	if err != nil {
		return nil, err
	}

	return result, nil
}
