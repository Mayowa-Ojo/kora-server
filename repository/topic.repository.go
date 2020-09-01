package repository

import (
	"github.com/Mayowa-Ojo/kora/config"
	"github.com/Mayowa-Ojo/kora/domain"
	"github.com/Mayowa-Ojo/kora/entity"
	"github.com/gofiber/fiber"
	"go.mongodb.org/mongo-driver/mongo"
)

// Topic - acts as the data access layer/repository for topics
type Topic struct {
	DB *mongo.Database
}

// NewTopicRepository -
func NewTopicRepository(conn *config.DBConn) domain.TopicRepository {
	return &Topic{conn.DB}
}

// GetAll -
func (t Topic) GetAll(ctx *fiber.Ctx) ([]entity.Topic, error) {
	return nil, nil
}
