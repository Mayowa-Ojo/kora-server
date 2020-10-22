package domain

import (
	"github.com/Mayowa-Ojo/kora/types"
	"github.com/gofiber/fiber"
)

// AuthService -
type AuthService interface {
	Login(ctx *fiber.Ctx) (types.GenericMap, error)
	Signup(ctx *fiber.Ctx) (types.GenericMap, error)
}
