package services

import (
	"strings"

	"github.com/Mayowa-Ojo/kora/constants"
	"github.com/Mayowa-Ojo/kora/domain"
	"github.com/Mayowa-Ojo/kora/entity"
	"github.com/Mayowa-Ojo/kora/utils"
	"github.com/fatih/structs"
	"github.com/gofiber/fiber"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// UserService - business logic layer/service for users
type UserService struct {
	userRepo  domain.UserRepository
	postRepo  domain.PostRepository
	topicRepo domain.TopicRepository
}

// NewUserService -
func NewUserService(
	u domain.UserRepository,
	p domain.PostRepository,
	t domain.TopicRepository,
) domain.UserService {
	return &UserService{
		u,
		p,
		t,
	}
}

// GetAll -
func (u *UserService) GetAll(ctx *fiber.Ctx) ([]entity.User, error) {
	return nil, nil
}

// GetOne -
func (u *UserService) GetOne(ctx *fiber.Ctx) (*entity.User, error) {
	userID := ctx.Params("id")
	userObjectID, err := primitive.ObjectIDFromHex(userID)
	filter := bson.D{{Key: "_id", Value: userObjectID}}
	opts := options.FindOne()

	user, err := u.userRepo.GetOne(ctx, filter, opts)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, constants.ErrNotFound
		}

		return nil, constants.ErrInternalServer
	}

	return user, nil
}

// UpdateContentViews - update views count when a post is requested
func (u *UserService) UpdateContentViews(ctx *fiber.Ctx) error {
	var filter bson.D
	postID := ctx.Params("postId")
	slug := ctx.Query("slug")
	username := ctx.Query("username")

	if postID != "" {
		postObjectID, err := primitive.ObjectIDFromHex(postID)
		if err != nil {
			return constants.ErrUnprocessableEntity
		}
		filter = bson.D{{Key: "_id", Value: postObjectID}}
	}

	if username != "" {
		filter = bson.D{
			{Key: "author.username", Value: username},
			{Key: "post_type", Value: "answer"},
		}
	}

	if slug != "" && username == "" {
		filter = bson.D{
			{Key: "slug", Value: slug},
			{Key: "post_type", Value: "question"},
		}
	}

	opts := options.FindOne()
	opts.SetProjection(bson.D{{Key: "author", Value: 1}})

	post, err := u.postRepo.GetOne(ctx, filter, opts)
	if err != nil {
		return constants.ErrNotFound
	}

	filter = bson.D{{Key: "_id", Value: post.Author.ID}}
	update := bson.D{{Key: "$inc", Value: bson.D{{Key: "views", Value: 1}}}}

	_, err = u.userRepo.UpdateOne(ctx, filter, update)
	if err != nil {
		return constants.ErrInternalServer
	}

	if username != "" {
		filter = bson.D{{Key: "_id", Value: post.ID}}
		update := bson.D{{Key: "$inc", Value: bson.D{{Key: "views", Value: 1}}}}

		_, err = u.postRepo.UpdateOne(ctx, filter, update)
		if err != nil {
			return constants.ErrInternalServer
		}
	}

	return nil
}

// UpdateProfile -
func (u *UserService) UpdateProfile(ctx *fiber.Ctx) (*entity.User, error) {
	var requestBody struct {
		Firstname  string `json:"firstname"`
		Lastname   string `json:"lastname"`
		About      string `json:"about"`
		Credential string `json:"credential"`
		Avatar     string `json:"avatar"`
	}

	userID, err := utils.GetJwtClaims(ctx, "userId")
	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, constants.ErrUnauthorized
	}

	if err := ctx.BodyParser(&requestBody); err != nil {
		return nil, constants.ErrUnprocessableEntity
	}

	mapRequestBody := structs.Map(&requestBody)

	filter := bson.D{{Key: "_id", Value: userObjectID}}
	update := bson.D{{Key: "$set", Value: bson.D{}}}

	for k, v := range mapRequestBody {
		if v != "" {
			update[0].Value = append(update[0].Value.(bson.D), bson.E{Key: strings.ToLower(k), Value: v})
		}
	}

	if _, err = u.userRepo.UpdateOne(ctx, filter, update); err != nil {
		return nil, constants.ErrInternalServer
	}

	filter = bson.D{{Key: "_id", Value: userObjectID}}
	opts := options.FindOne()
	user, err := u.userRepo.GetOne(ctx, filter, opts)

	if err != nil {
		return nil, constants.ErrInternalServer
	}

	return user, nil
}

// GetFollowersForUser -
func (u *UserService) GetFollowersForUser(ctx *fiber.Ctx) ([]entity.User, error) {
	user, err := utils.GetUserFromAuthHeader(ctx, u.userRepo)
	if err != nil {
		return nil, constants.ErrUnauthorized
	}

	filter := bson.D{{Key: "_id", Value: bson.D{{Key: "$in", Value: user.Followers}}}}
	opts := options.Find()

	followers, err := u.userRepo.GetMany(ctx, filter, opts)
	if err != nil {
		return nil, constants.ErrInternalServer
	}

	return followers, nil
}

// GetFollowingForUser -
func (u *UserService) GetFollowingForUser(ctx *fiber.Ctx) ([]entity.User, error) {
	user, err := utils.GetUserFromAuthHeader(ctx, u.userRepo)
	if err != nil {
		return nil, constants.ErrUnauthorized
	}

	filter := bson.D{{Key: "_id", Value: bson.D{{Key: "$in", Value: user.Following}}}}
	opts := options.Find()

	followers, err := u.userRepo.GetMany(ctx, filter, opts)
	if err != nil {
		return nil, constants.ErrInternalServer
	}

	return followers, nil
}

// GetPostsForUser -
func (u *UserService) GetPostsForUser(ctx *fiber.Ctx) ([]entity.Post, error) {
	postType := ctx.Query("postType")

	if postType == "" {
		return nil, constants.ErrUnprocessableEntity
	}

	user, err := utils.GetUserFromAuthHeader(ctx, u.userRepo)
	if err != nil {
		return nil, constants.ErrUnauthorized
	}

	var filter bson.D
	switch postType {
	case "question":
		filter = bson.D{{Key: "_id", Value: bson.D{{Key: "$in", Value: user.Questions}}}}
	case "answer":
		filter = bson.D{{Key: "_id", Value: bson.D{{Key: "$in", Value: user.Answers}}}}
	case "post":
		filter = bson.D{{Key: "_id", Value: bson.D{{Key: "$in", Value: user.Posts}}}}
	}
	opts := options.Find()

	followers, err := u.postRepo.GetMany(ctx, filter, opts)
	if err != nil {
		return nil, constants.ErrInternalServer
	}

	return followers, nil
}

// GetSpacesForUser -
func (u *UserService) GetSpacesForUser(ctx *fiber.Ctx) ([]entity.Post, error) {
	return nil, nil
}

// GetKnowledgeForUser -
func (u *UserService) GetKnowledgeForUser(ctx *fiber.Ctx) ([]entity.Topic, error) {
	user, err := utils.GetUserFromAuthHeader(ctx, u.userRepo)
	if err != nil {
		return nil, constants.ErrUnauthorized
	}

	filter := bson.D{{Key: "_id", Value: bson.D{{Key: "$in", Value: user.Knowledge}}}}
	opts := options.Find()

	topics, err := u.topicRepo.GetMany(ctx, filter, opts)
	if err != nil {
		return nil, constants.ErrInternalServer
	}

	return topics, nil
}

// FollowUser -
func (u *UserService) FollowUser(ctx *fiber.Ctx) error {
	user, err := utils.GetUserFromAuthHeader(ctx, u.userRepo)
	if err != nil {
		return constants.ErrUnauthorized
	}

	followedUserID := ctx.Query("userId")
	followedUserObjectID, err := primitive.ObjectIDFromHex(followedUserID)
	if err != nil {
		return constants.ErrInternalServer
	}

	filter := bson.D{{Key: "_id", Value: followedUserObjectID}}
	update := bson.D{{Key: "$push", Value: bson.D{{Key: "followers", Value: user.ID}}}}
	if _, err := u.userRepo.UpdateOne(ctx, filter, update); err != nil {
		return constants.ErrInternalServer
	}

	filter = bson.D{{Key: "_id", Value: user.ID}}
	update = bson.D{{Key: "$push", Value: bson.D{{Key: "following", Value: followedUserObjectID}}}}
	if _, err := u.userRepo.UpdateOne(ctx, filter, update); err != nil {
		return constants.ErrInternalServer
	}

	return nil
}

// UnfollowUser -
func (u *UserService) UnfollowUser(ctx *fiber.Ctx) error {
	user, err := utils.GetUserFromAuthHeader(ctx, u.userRepo)
	if err != nil {
		return constants.ErrUnauthorized
	}

	followedUserID := ctx.Query("userId")
	followedUserObjectID, err := primitive.ObjectIDFromHex(followedUserID)
	if err != nil {
		return constants.ErrInternalServer
	}

	filter := bson.D{{Key: "_id", Value: followedUserObjectID}}
	update := bson.D{{Key: "$pull", Value: bson.D{{Key: "followers", Value: user.ID}}}}
	if _, err := u.userRepo.UpdateOne(ctx, filter, update); err != nil {
		return constants.ErrInternalServer
	}

	filter = bson.D{{Key: "_id", Value: user.ID}}
	update = bson.D{{Key: "$pull", Value: bson.D{{Key: "following", Value: followedUserObjectID}}}}
	if _, err := u.userRepo.UpdateOne(ctx, filter, update); err != nil {
		return constants.ErrInternalServer
	}

	return nil
}

// SetPinnedPost -
func (u *UserService) SetPinnedPost(ctx *fiber.Ctx) error {
	postID := ctx.Query("postId")
	postObjectID, err := primitive.ObjectIDFromHex(postID)
	if err != nil {
		return constants.ErrUnauthorized
	}

	userID, err := utils.GetJwtClaims(ctx, "userId")
	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return constants.ErrUnauthorized
	}

	filter := bson.D{{Key: "_id", Value: userObjectID}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "pinned_post", Value: postObjectID}}}}
	if _, err := u.userRepo.UpdateOne(ctx, filter, update); err != nil {
		return constants.ErrInternalServer
	}

	return nil
}

// UnsetPinnedPost -
func (u *UserService) UnsetPinnedPost(ctx *fiber.Ctx) error {
	userID, err := utils.GetJwtClaims(ctx, "userId")
	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return constants.ErrUnauthorized
	}

	filter := bson.D{{Key: "_id", Value: userObjectID}}
	update := bson.D{{Key: "$unset", Value: bson.D{{Key: "pinned_post", Value: ""}}}}
	if _, err := u.userRepo.UpdateOne(ctx, filter, update); err != nil {
		return constants.ErrInternalServer
	}

	return nil
}
