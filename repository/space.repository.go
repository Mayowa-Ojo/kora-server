package repository

import (
	"github.com/Mayowa-Ojo/kora-server/config"
	"github.com/Mayowa-Ojo/kora-server/domain"
	"github.com/Mayowa-Ojo/kora-server/entity"
	"github.com/Mayowa-Ojo/kora-server/types"
	"github.com/gofiber/fiber"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// SpaceRepository - acts as the data access layer/repository for spaces
type SpaceRepository struct {
	DB *mongo.Database
}

// NewSpaceRepository -
func NewSpaceRepository(conn *config.DBConn) domain.SpaceRepository {
	return &SpaceRepository{conn.DB}
}

// GetAll -
func (s *SpaceRepository) GetAll(ctx *fiber.Ctx) ([]entity.Space, error) {
	col := s.DB.Collection("spaces")
	spaces := make([]entity.Space, 0)
	filter := bson.D{{}}

	cur, err := col.Find(ctx.Fasthttp, filter)
	if err != nil {
		return nil, err
	}

	if err := cur.All(ctx.Fasthttp, &spaces); err != nil {
		return nil, err
	}

	return spaces, nil

}

// GetOne -
func (s *SpaceRepository) GetOne(ctx *fiber.Ctx, filter types.Any, opts *options.FindOneOptions) (*entity.Space, error) {
	col := s.DB.Collection("spaces")
	space := new(entity.Space)

	result := col.FindOne(ctx.Fasthttp, filter, opts)
	if err := result.Decode(&space); err != nil {
		return nil, err
	}

	return space, nil
}

// Create -
func (s *SpaceRepository) Create(ctx *fiber.Ctx, space *entity.Space) (*mongo.InsertOneResult, error) {
	col := s.DB.Collection("spaces")

	result, err := col.InsertOne(ctx.Fasthttp, space)

	if err != nil {
		return nil, err
	}

	return result, nil
}

// GetMany -
func (s *SpaceRepository) GetMany(ctx *fiber.Ctx, filter types.Any, opts *options.FindOptions) ([]entity.Space, error) {
	col := s.DB.Collection("spaces")
	spaces := make([]entity.Space, 0)

	cur, err := col.Find(ctx.Fasthttp, filter, opts)
	if err != nil {
		return nil, err
	}

	if err := cur.All(ctx.Fasthttp, &spaces); err != nil {
		return nil, err
	}

	return spaces, nil
}

// UpdateOne -
func (s *SpaceRepository) UpdateOne(ctx *fiber.Ctx, filter, update types.Any) (*mongo.UpdateResult, error) {
	col := s.DB.Collection("spaces")

	result, err := col.UpdateOne(ctx.Fasthttp, filter, update)

	if err != nil {
		return nil, err
	}

	return result, nil
}

// DeleteOne -
func (s *SpaceRepository) DeleteOne(ctx *fiber.Ctx, filter types.Any) (*mongo.DeleteResult, error) {
	col := s.DB.Collection("spaces")

	result, err := col.DeleteOne(ctx.Fasthttp, filter)

	if err != nil {
		return nil, err
	}

	return result, nil
}
