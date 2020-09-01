package domain

import (
	"github.com/Mayowa-Ojo/kora/entity"
	"github.com/gofiber/fiber"
)

// TopicService -
type TopicService interface {
	GetAll(ctx *fiber.Ctx) ([]entity.Topic, error)
}

// TopicRepository -
type TopicRepository interface {
	GetAll(ctx *fiber.Ctx) ([]entity.Topic, error)
}
