package domain

import (
	"github.com/Mayowa-Ojo/kora/entity"
	"github.com/gofiber/fiber"
)

// UserService -
type UserService interface {
	GetAll(ctx *fiber.Ctx) ([]entity.User, error)
}

// UserRepository -
type UserRepository interface {
	GetAll(ctx *fiber.Ctx) ([]entity.User, error)
}
