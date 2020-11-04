package domain

import (
	"github.com/Mayowa-Ojo/kora/types"
	"github.com/gofiber/fiber"
)

// ImageService -
type ImageService interface {
	UploadImage(ctx *fiber.Ctx) (types.GenericMap, error)
	GetImage(ctx *fiber.Ctx) error
	DeleteImage(ctx *fiber.Ctx) error
}
