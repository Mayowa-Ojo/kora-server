package domain

import (
	"github.com/Mayowa-Ojo/kora/entity"
	"github.com/Mayowa-Ojo/kora/types"
	"github.com/gofiber/fiber"
	mg "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// UserService -
type UserService interface {
	GetAll(ctx *fiber.Ctx) ([]entity.User, error)
}

// UserRepository -
type UserRepository interface {
	GetAll(ctx *fiber.Ctx) ([]entity.User, error)
	GetOne(ctx *fiber.Ctx, filter types.Any, opts *options.FindOneOptions) (*entity.User, error)
	Create(ctx *fiber.Ctx, user entity.User) (*mg.InsertOneResult, error)
}
