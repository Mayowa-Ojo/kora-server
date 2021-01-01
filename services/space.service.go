package services

import (
	"fmt"

	"github.com/Mayowa-Ojo/kora/constants"
	"github.com/Mayowa-Ojo/kora/domain"
	"github.com/Mayowa-Ojo/kora/entity"
	"github.com/Mayowa-Ojo/kora/types"
	"github.com/Mayowa-Ojo/kora/utils"
	snake "github.com/Mayowa-Ojo/rattle-snake"
	"github.com/fatih/structs"
	"github.com/gofiber/fiber"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// SpaceService -
type SpaceService struct {
	spaceRepo domain.SpaceRepository
	postRepo  domain.PostRepository
	userRepo  domain.UserRepository
}

// NewSpaceService -
func NewSpaceService(
	s domain.SpaceRepository,
	p domain.PostRepository,
	u domain.UserRepository,
) domain.SpaceService {
	return &SpaceService{
		s,
		p,
		u,
	}
}

// GetAll -
func (s *SpaceService) GetAll(ctx *fiber.Ctx) ([]entity.Space, error) {
	spaces, err := s.spaceRepo.GetAll(ctx)

	if err != nil {
		return nil, constants.ErrInternalServer
	}

	return spaces, nil
}

// GetOne -
func (s *SpaceService) GetOne(ctx *fiber.Ctx) (*entity.Space, error) {
	spaceID := ctx.Params("id")
	spaceObjectID, err := primitive.ObjectIDFromHex(spaceID)
	filter := bson.D{{Key: "_id", Value: spaceObjectID}}
	opts := options.FindOne()

	space, err := s.spaceRepo.GetOne(ctx, filter, opts)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, constants.ErrNotFound
		}

		return nil, constants.ErrInternalServer
	}

	return space, nil
}

// GetBySlug -
func (s *SpaceService) GetBySlug(ctx *fiber.Ctx) (*entity.Space, error) {
	slug := ctx.Query("q")

	filter := bson.D{{Key: "slug", Value: slug}}
	opts := options.FindOne()

	space, err := s.spaceRepo.GetOne(ctx, filter, opts)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, constants.ErrNotFound
		}

		return nil, constants.ErrInternalServer
	}

	return space, nil
}

// Create -
func (s *SpaceService) Create(ctx *fiber.Ctx) (*entity.Space, error) {
	var requestBody struct {
		Name  string `json:"name"`
		About string `json:"about"`
	}

	err := ctx.BodyParser(&requestBody)
	if err != nil {
		return nil, constants.ErrUnprocessableEntity
	}

	userID, err := utils.GetJwtClaims(ctx, "userId")
	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, constants.ErrUnauthorized
	}

	slug := utils.GenerateSlug(requestBody.Name)

	instance := &entity.Space{
		Name:  requestBody.Name,
		About: requestBody.About,
		Slug:  slug,
	}

	if err := instance.Validate(); err != nil {
		return nil, constants.ErrUnprocessableEntity
	}

	instance.SetDefaultValues()

	insertResult, err := s.spaceRepo.Create(ctx, instance)
	if err != nil {
		return nil, constants.ErrInternalServer
	}

	filter := bson.D{{Key: "_id", Value: insertResult.InsertedID}}
	update := bson.D{{Key: "$addToSet", Value: bson.D{{Key: "admins", Value: userObjectID}}}}
	if _, err := s.spaceRepo.UpdateOne(ctx, filter, update); err != nil {
		return nil, constants.ErrInternalServer
	}

	opts := options.FindOne()
	space, err := s.spaceRepo.GetOne(ctx, filter, opts)
	if err != nil {
		return nil, constants.ErrInternalServer
	}

	return space, nil
}

// GetPostsForSpace -
// [E.g] /api/v1/spaces/:id/posts?postType=question
func (s *SpaceService) GetPostsForSpace(ctx *fiber.Ctx) ([]entity.Post, error) {
	postType := ctx.Query("postType")
	if postType != "all" && postType != "question" {
		return nil, constants.ErrInvalidCredentials
	}

	spaceID := ctx.Params("id")
	spaceObjectID, err := primitive.ObjectIDFromHex(spaceID)
	if err != nil {
		return nil, constants.ErrUnprocessableEntity
	}

	filter := bson.D{{Key: "_id", Value: spaceObjectID}}
	opts := options.FindOne()
	space, err := s.spaceRepo.GetOne(ctx, filter, opts)
	if err != nil {
		return nil, constants.ErrNotFound
	}

	switch postType {
	case "all":
		filter = bson.D{{Key: "_id", Value: bson.D{{Key: "$in", Value: space.Posts}}}}
	case "questions":
		filter = bson.D{{
			Key: "$and",
			Value: bson.A{
				bson.D{{Key: "_id", Value: bson.D{{Key: "$in", Value: space.Posts}}}},
				bson.D{{Key: "postType", Value: "question"}},
			},
		}}
	}

	postOpts := options.Find()
	posts, err := s.postRepo.GetMany(ctx, filter, postOpts)
	if err != nil {
		return nil, constants.ErrInternalServer
	}

	return posts, nil
}

// GetMembersForSpace - /api/v1/spaces/:id/people
func (s *SpaceService) GetMembersForSpace(ctx *fiber.Ctx) (types.GenericMap, error) {
	spaceID := ctx.Params("id")
	spaceObjectID, err := primitive.ObjectIDFromHex(spaceID)
	if err != nil {
		return nil, constants.ErrUnprocessableEntity
	}

	filter := bson.D{{Key: "_id", Value: spaceObjectID}}
	spaceOpts := options.FindOne()
	space, err := s.spaceRepo.GetOne(ctx, filter, spaceOpts)
	if err != nil {
		return nil, constants.ErrNotFound
	}

	filter = bson.D{{Key: "_id", Value: bson.D{{Key: "$in", Value: space.Followers}}}}
	opts := options.Find()
	followers, err := s.userRepo.GetMany(ctx, filter, opts)

	filter = bson.D{{Key: "_id", Value: bson.D{{Key: "$in", Value: space.Contributors}}}}
	contributors, err := s.userRepo.GetMany(ctx, filter, opts)

	filter = bson.D{{Key: "_id", Value: bson.D{{Key: "$in", Value: space.Moderators}}}}
	moderators, err := s.userRepo.GetMany(ctx, filter, opts)

	filter = bson.D{{Key: "_id", Value: bson.D{{Key: "$in", Value: space.Admins}}}}
	admins, err := s.userRepo.GetMany(ctx, filter, opts)

	if err != nil {
		return nil, constants.ErrInternalServer
	}

	result := types.GenericMap{
		"followers":    followers,
		"contributors": contributors,
		"moderators":   moderators,
		"admins":       admins,
	}

	return result, nil
}

// GetSuggestedSpaces - /api/v1/spaces/suggestions
func (s *SpaceService) GetSuggestedSpaces(ctx *fiber.Ctx) ([]entity.Space, error) {
	// should implement better suggestion logic
	userID, err := utils.GetJwtClaims(ctx, "userId")
	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, constants.ErrUnauthorized
	}

	filter := bson.D{{Key: "followers", Value: bson.D{{Key: "$nin", Value: []primitive.ObjectID{userObjectID}}}}}
	opts := options.Find()
	opts.SetLimit(20)

	spaces, err := s.spaceRepo.GetMany(ctx, filter, opts)
	if err != nil {
		fmt.Println(err)
		return nil, constants.ErrInternalServer
	}

	return spaces, nil
}

// UpdateProfileByAdmin -
func (s *SpaceService) UpdateProfileByAdmin(ctx *fiber.Ctx) (*entity.Space, error) {
	var requestBody struct {
		Name       string `json:"name"`
		About      string `json:"about"`
		Details    string `json:"details"`
		Icon       string `json:"icon"`
		CoverPhoto string `json:"coverPhoto"`
	}

	spaceID := ctx.Params("id")
	spaceObjectID, err := primitive.ObjectIDFromHex(spaceID)
	if err != nil {
		return nil, constants.ErrUnprocessableEntity
	}

	if err := ctx.BodyParser(&requestBody); err != nil {
		return nil, constants.ErrUnprocessableEntity
	}

	mapRequestBody := structs.Map(&requestBody)

	filter := bson.D{{Key: "_id", Value: spaceObjectID}}
	update := bson.D{{Key: "$set", Value: bson.D{}}}

	for k, v := range mapRequestBody {
		if v != "" {
			update[0].Value = append(update[0].Value.(bson.D), bson.E{Key: snake.ToSnakeCase(k), Value: v})
		}
	}

	if _, err = s.spaceRepo.UpdateOne(ctx, filter, update); err != nil {
		return nil, constants.ErrInternalServer
	}

	filter = bson.D{{Key: "_id", Value: spaceObjectID}}
	opts := options.FindOne()
	space, err := s.spaceRepo.GetOne(ctx, filter, opts)

	if err != nil {
		return nil, constants.ErrInternalServer
	}

	return space, nil
}

// DeleteSpaceByAdmin -
func (s *SpaceService) DeleteSpaceByAdmin(ctx *fiber.Ctx) error {
	spaceID := ctx.Params("id")
	spaceObjectID, err := primitive.ObjectIDFromHex(spaceID)
	if err != nil {
		return constants.ErrUnprocessableEntity
	}

	filter := bson.D{{Key: "_id", Value: spaceObjectID}}
	opts := options.FindOne()
	space, err := s.spaceRepo.GetOne(ctx, filter, opts)
	if err != nil {
		return constants.ErrNotFound
	}

	userID, err := utils.GetJwtClaims(ctx, "userId")
	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return constants.ErrUnauthorized
	}

	for _, v := range space.Admins {
		if v == userObjectID {

			_, err := s.spaceRepo.DeleteOne(ctx, filter)
			if err != nil {
				return constants.ErrInternalServer
			}

			return nil
		}
	}

	return constants.ErrForbidden
}

// FollowSpace -
func (s *SpaceService) FollowSpace(ctx *fiber.Ctx) (*entity.User, error) {
	spaceID := ctx.Params("id")
	spaceObjectID, err := primitive.ObjectIDFromHex(spaceID)
	if err != nil {
		return nil, constants.ErrUnprocessableEntity
	}

	userID, err := utils.GetJwtClaims(ctx, "userId")
	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, constants.ErrUnauthorized
	}

	filter := bson.D{{Key: "_id", Value: spaceObjectID}}
	update := bson.D{{Key: "$addToSet", Value: bson.D{{Key: "followers", Value: userObjectID}}}}
	_, err = s.spaceRepo.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, constants.ErrInternalServer
	}

	filter = bson.D{{Key: "_id", Value: userObjectID}}
	update = bson.D{{Key: "$addToSet", Value: bson.D{{Key: "spaces", Value: spaceObjectID}}}}
	_, err = s.userRepo.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, constants.ErrInternalServer
	}

	filter = bson.D{{Key: "_id", Value: userObjectID}}
	opts := options.FindOne()

	user, err := s.userRepo.GetOne(ctx, filter, opts)
	if err != nil {
		return nil, constants.ErrInternalServer
	}

	return user, nil
}

// UnfollowSpace -
func (s *SpaceService) UnfollowSpace(ctx *fiber.Ctx) (*entity.User, error) {
	spaceID := ctx.Params("id")
	spaceObjectID, err := primitive.ObjectIDFromHex(spaceID)
	if err != nil {
		return nil, constants.ErrUnprocessableEntity
	}

	userID, err := utils.GetJwtClaims(ctx, "userId")
	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, constants.ErrUnauthorized
	}

	filter := bson.D{{Key: "_id", Value: spaceObjectID}}
	update := bson.D{{Key: "$pull", Value: bson.D{{Key: "followers", Value: userObjectID}}}}
	_, err = s.spaceRepo.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, constants.ErrInternalServer
	}

	filter = bson.D{{Key: "_id", Value: userObjectID}}
	update = bson.D{{Key: "$pull", Value: bson.D{{Key: "spaces", Value: spaceObjectID}}}}
	_, err = s.userRepo.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, constants.ErrInternalServer
	}

	filter = bson.D{{Key: "_id", Value: userObjectID}}
	opts := options.FindOne()

	user, err := s.userRepo.GetOne(ctx, filter, opts)
	if err != nil {
		return nil, constants.ErrInternalServer
	}

	return user, nil
}

// SetPinnedPost -
func (s *SpaceService) SetPinnedPost(ctx *fiber.Ctx) error {
	postID := ctx.Query("postId")
	postObjectID, err := primitive.ObjectIDFromHex(postID)
	if err != nil {
		return constants.ErrUnauthorized
	}

	spaceID := ctx.Params("id")
	spaceObjectID, err := primitive.ObjectIDFromHex(spaceID)
	if err != nil {
		return constants.ErrUnprocessableEntity
	}

	filter := bson.D{{Key: "_id", Value: spaceObjectID}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "pinned_post", Value: postObjectID}}}}
	if _, err := s.spaceRepo.UpdateOne(ctx, filter, update); err != nil {
		return constants.ErrInternalServer
	}

	return nil
}

// UnsetPinnedPost -
func (s *SpaceService) UnsetPinnedPost(ctx *fiber.Ctx) error {
	postID := ctx.Query("postId")
	postObjectID, err := primitive.ObjectIDFromHex(postID)
	if err != nil {
		return constants.ErrUnauthorized
	}

	spaceID := ctx.Params("id")
	spaceObjectID, err := primitive.ObjectIDFromHex(spaceID)
	if err != nil {
		return constants.ErrUnprocessableEntity
	}

	filter := bson.D{{Key: "_id", Value: spaceObjectID}}
	update := bson.D{{Key: "$unset", Value: bson.D{{Key: "pinned_post", Value: postObjectID}}}}
	if _, err := s.spaceRepo.UpdateOne(ctx, filter, update); err != nil {
		return constants.ErrInternalServer
	}

	return nil
}
