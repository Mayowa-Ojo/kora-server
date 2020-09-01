package repository

import (
	"github.com/Mayowa-Ojo/kora/config"
	"github.com/Mayowa-Ojo/kora/domain"
	"github.com/Mayowa-Ojo/kora/entity"
	"github.com/gofiber/fiber"
	"go.mongodb.org/mongo-driver/mongo"
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
	return nil, nil
}
