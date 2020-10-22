package repository

import (
	"github.com/Mayowa-Ojo/kora/config"
	"github.com/Mayowa-Ojo/kora/domain"
	"github.com/Mayowa-Ojo/kora/entity"
	"github.com/Mayowa-Ojo/kora/types"
	"github.com/gofiber/fiber"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	mg "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// User - acts as the data access layer/repository for users
type User struct {
	DB *mongo.Database
}

// NewUserRepository -
func NewUserRepository(conn *config.DBConn) domain.UserRepository {
	return &User{conn.DB}
}

// GetAll -
func (u User) GetAll(ctx *fiber.Ctx) ([]entity.User, error) {
	users := make([]entity.User, 0)
	filter := bson.D{{}}
	col := u.DB.Collection("users")

	c, err := col.Find(ctx.Fasthttp, filter)
	if err != nil {
		return nil, err
	}

	if err := c.All(ctx.Fasthttp, &users); err != nil {
		return nil, err
	}

	return users, nil
}

// GetOne -
func (u User) GetOne(ctx *fiber.Ctx, filter types.Any, opts *options.FindOneOptions) (*entity.User, error) {
	c := u.DB.Collection("users")
	user := new(entity.User)

	result := c.FindOne(ctx.Fasthttp, filter, opts)

	if err := result.Decode(&user); err != nil {
		return nil, err
	}

	return user, nil
}

// Create -
func (u User) Create(ctx *fiber.Ctx, user entity.User) (*mg.InsertOneResult, error) {
	c := u.DB.Collection("users")

	result, err := c.InsertOne(ctx.Fasthttp, user)

	if err != nil {
		return nil, err
	}

	return result, nil
}
