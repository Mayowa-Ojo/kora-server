package services

import (
	"fmt"
	"strings"
	"time"

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
	userRepo       domain.UserRepository
	postRepo       domain.PostRepository
	topicRepo      domain.TopicRepository
	spaceRepo      domain.SpaceRepository
	sharedPostRepo domain.SharedPostRepository
}

// NewUserService -
func NewUserService(
	u domain.UserRepository,
	p domain.PostRepository,
	t domain.TopicRepository,
	s domain.SpaceRepository,
	sp domain.SharedPostRepository,
) domain.UserService {
	return &UserService{
		u,
		p,
		t,
		s,
		sp,
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
	opts.SetProjection(bson.D{{Key: "hash", Value: 0}})

	user, err := u.userRepo.GetOne(ctx, filter, opts)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, constants.ErrNotFound
		}

		return nil, constants.ErrInternalServer
	}

	time.Sleep(time.Second * 2)

	return user, nil
}

// GetUserProfile -
func (u *UserService) GetUserProfile(ctx *fiber.Ctx) (*entity.User, error) {
	username := ctx.Query("q")

	filter := bson.D{{Key: "username", Value: username}}
	opts := options.FindOne()
	opts.SetProjection(bson.D{{Key: "hash", Value: 0}})

	user, err := u.userRepo.GetOne(ctx, filter, opts)
	if err != nil {
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
			{Key: "slug", Value: slug},
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
		filter := bson.D{{Key: "_id", Value: post.ID}}
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
	type Credentials struct {
		Profile    string `json:"profile"`
		Employment string `json:"employment"`
		Education  string `json:"education"`
		Location   string `json:"location"`
		Language   string `json:"language"`
		Space      string `json:"space"`
		Topic      string `json:"topic"`
	}
	var requestBody struct {
		Firstname   string      `json:"firstname"`
		Lastname    string      `json:"lastname"`
		About       string      `json:"about"`
		Credentials Credentials `json:"credentials"`
		Avatar      string      `json:"avatar"`
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
		if _, ok := v.(map[string]interface{}); ok { // assert a nested struct field and handle its update

			for kk, vv := range v.(map[string]interface{}) {
				if vv != "" {
					update[0].Value = append(
						update[0].Value.(bson.D), bson.E{Key: strings.ToLower(fmt.Sprintf("%s.%s", k, kk)), Value: vv},
					)
				}
			}

			continue
		}

		if v != "" {
			update[0].Value = append(update[0].Value.(bson.D), bson.E{Key: strings.ToLower(k), Value: v})
		}
	}

	if _, err = u.userRepo.UpdateOne(ctx, filter, update); err != nil {
		return nil, constants.ErrInternalServer
	}

	filter = bson.D{{Key: "_id", Value: userObjectID}}
	opts := options.FindOne()
	opts.SetProjection(bson.D{{Key: "hash", Value: 0}})

	user, err := u.userRepo.GetOne(ctx, filter, opts)
	if err != nil {
		return nil, constants.ErrInternalServer
	}

	return user, nil
}

// UpdateUserKnowledge -
func (u *UserService) UpdateUserKnowledge(ctx *fiber.Ctx) ([]entity.Topic, error) {
	var requestBody struct {
		Knowledge []string `json:"knowledge"`
	}

	userID, err := utils.GetJwtClaims(ctx, "userId")
	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, constants.ErrUnauthorized
	}

	if err := ctx.BodyParser(&requestBody); err != nil {
		return nil, constants.ErrUnprocessableEntity
	}

	for _, v := range requestBody.Knowledge {
		filter := bson.D{{Key: "name", Value: v}}
		_, err := u.topicRepo.GetOne(ctx, filter)

		if err != nil && err != mongo.ErrNoDocuments {
			return nil, constants.ErrInternalServer
		}

		if err == mongo.ErrNoDocuments {
			// create topic
			instance := &entity.Topic{
				Name: v,
			}

			if err := instance.Validate(); err != nil {
				return nil, constants.ErrUnprocessableEntity
			}

			instance.SetDefaultValues()

			_, err := u.topicRepo.Create(ctx, instance)
			if err != nil {
				return nil, constants.ErrInternalServer
			}
		}
	}

	filter := bson.D{{Key: "name", Value: bson.D{{Key: "$in", Value: requestBody.Knowledge}}}}
	opts := options.Find()

	topics, err := u.topicRepo.GetMany(ctx, filter, opts)
	if err != nil {
		return nil, constants.ErrInternalServer
	}

	var topicIDs []primitive.ObjectID
	for _, v := range topics {
		topicIDs = append(topicIDs, v.ID)
	}

	filter = bson.D{{Key: "_id", Value: userObjectID}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "knowledge", Value: topicIDs}}}}

	if _, err := u.userRepo.UpdateOne(ctx, filter, update); err != nil {
		return nil, constants.ErrInternalServer
	}

	return topics, nil
}

// GetFollowersForUser -
func (u *UserService) GetFollowersForUser(ctx *fiber.Ctx) ([]entity.User, error) {
	userID := ctx.Params("id")
	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, constants.ErrUnprocessableEntity
	}

	filter := bson.D{{Key: "_id", Value: userObjectID}}
	opts := options.FindOne()

	user, err := u.userRepo.GetOne(ctx, filter, opts)
	if err != nil {
		return nil, constants.ErrInternalServer
	}

	filter = bson.D{{Key: "_id", Value: bson.D{{Key: "$in", Value: user.Followers}}}}
	findOpts := options.Find()

	followers, err := u.userRepo.GetMany(ctx, filter, findOpts)
	if err != nil {
		return nil, constants.ErrInternalServer
	}

	return followers, nil
}

// GetFollowingForUser -
func (u *UserService) GetFollowingForUser(ctx *fiber.Ctx) ([]entity.User, error) {
	userID := ctx.Params("id")
	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, constants.ErrUnprocessableEntity
	}

	filter := bson.D{{Key: "_id", Value: userObjectID}}
	opts := options.FindOne()

	user, err := u.userRepo.GetOne(ctx, filter, opts)
	if err != nil {
		return nil, constants.ErrInternalServer
	}

	filter = bson.D{{Key: "_id", Value: bson.D{{Key: "$in", Value: user.Following}}}}
	findOpts := options.Find()

	following, err := u.userRepo.GetMany(ctx, filter, findOpts)
	if err != nil {
		return nil, constants.ErrInternalServer
	}

	return following, nil
}

// GetPostsForUser -
func (u *UserService) GetPostsForUser(ctx *fiber.Ctx) ([]entity.Post, error) {
	var filter bson.D

	postType := ctx.Query("postType")
	userID := ctx.Params("id")
	userObjectID, err := primitive.ObjectIDFromHex(userID)

	if err != nil {
		return nil, constants.ErrUnprocessableEntity
	}

	if postType == "" {
		return nil, constants.ErrUnprocessableEntity
	}

	switch postType {
	case "question":
		filter = bson.D{{Key: "post_type", Value: "question"}, {Key: "author._id", Value: userObjectID}}
	case "answer":
		filter = bson.D{{Key: "post_type", Value: "answer"}, {Key: "author._id", Value: userObjectID}}
	case "post":
		filter = bson.D{{Key: "post_type", Value: "post"}, {Key: "author._id", Value: userObjectID}}
	default:
		return nil, constants.ErrUnprocessableEntity
	}

	opts := options.Find()
	posts, err := u.postRepo.GetMany(ctx, filter, opts)
	if err != nil {
		return nil, constants.ErrInternalServer
	}

	return posts, nil
}

// GetSharedPostsForUser -
func (u *UserService) GetSharedPostsForUser(ctx *fiber.Ctx) ([]entity.SharedPost, error) {
	userID := ctx.Params("id")
	userObjectID, err := primitive.ObjectIDFromHex(userID)

	if err != nil {
		return nil, constants.ErrUnprocessableEntity
	}

	filter := bson.D{{Key: "author._id", Value: userObjectID}}
	opts := options.Find()

	sharedPosts, err := u.sharedPostRepo.GetMany(ctx, filter, opts)
	if err != nil {
		return nil, constants.ErrInternalServer
	}

	return sharedPosts, nil
}

// GetSpacesForUser -
func (u *UserService) GetSpacesForUser(ctx *fiber.Ctx) ([]entity.Space, error) {
	userID := ctx.Params("id")
	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, constants.ErrUnprocessableEntity
	}

	filter := bson.D{{Key: "followers", Value: bson.D{
		{Key: "$elemMatch", Value: bson.D{{Key: "$eq", Value: userObjectID}}},
	}}}
	opts := options.Find()

	spaces, err := u.spaceRepo.GetMany(ctx, filter, opts)
	if err != nil {
		return nil, constants.ErrInternalServer
	}

	return spaces, nil
}

// GetKnowledgeForUser -
func (u *UserService) GetKnowledgeForUser(ctx *fiber.Ctx) ([]entity.Topic, error) {
	userID := ctx.Params("id")
	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, constants.ErrUnprocessableEntity
	}

	filter := bson.D{{Key: "_id", Value: userObjectID}}
	opts := options.FindOne()

	user, err := u.userRepo.GetOne(ctx, filter, opts)
	if err != nil {
		return nil, constants.ErrInternalServer
	}

	filter = bson.D{{Key: "_id", Value: bson.D{{Key: "$in", Value: user.Knowledge}}}}
	findOpts := options.Find()

	topics, err := u.topicRepo.GetMany(ctx, filter, findOpts)
	if err != nil {
		return nil, constants.ErrInternalServer
	}

	return topics, nil
}

// FollowUser -
func (u *UserService) FollowUser(ctx *fiber.Ctx) (*entity.User, error) {
	userID, err := utils.GetJwtClaims(ctx, "userId")
	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, constants.ErrUnauthorized
	}

	followedUserID := ctx.Params("id")
	followedUserObjectID, err := primitive.ObjectIDFromHex(followedUserID)
	if err != nil {
		return nil, constants.ErrInternalServer
	}

	filter := bson.D{{Key: "_id", Value: followedUserObjectID}}
	update := bson.D{{Key: "$addToSet", Value: bson.D{{Key: "followers", Value: userObjectID}}}}
	if _, err := u.userRepo.UpdateOne(ctx, filter, update); err != nil {
		return nil, constants.ErrInternalServer
	}

	filter = bson.D{{Key: "_id", Value: userObjectID}}
	update = bson.D{{Key: "$addToSet", Value: bson.D{{Key: "following", Value: followedUserObjectID}}}}
	if _, err := u.userRepo.UpdateOne(ctx, filter, update); err != nil {
		return nil, constants.ErrInternalServer
	}

	filter = bson.D{{Key: "_id", Value: userObjectID}}
	opts := options.FindOne()
	opts.SetProjection(bson.D{{Key: "hash", Value: 0}})

	user, err := u.userRepo.GetOne(ctx, filter, opts)
	if err != nil {
		return nil, constants.ErrInternalServer
	}

	return user, nil
}

// UnfollowUser -
func (u *UserService) UnfollowUser(ctx *fiber.Ctx) (*entity.User, error) {
	userID, err := utils.GetJwtClaims(ctx, "userId")
	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, constants.ErrUnauthorized
	}

	followedUserID := ctx.Params("id")
	followedUserObjectID, err := primitive.ObjectIDFromHex(followedUserID)
	if err != nil {
		return nil, constants.ErrInternalServer
	}

	filter := bson.D{{Key: "_id", Value: followedUserObjectID}}
	update := bson.D{{Key: "$pull", Value: bson.D{{Key: "followers", Value: userObjectID}}}}
	if _, err := u.userRepo.UpdateOne(ctx, filter, update); err != nil {
		return nil, constants.ErrInternalServer
	}

	filter = bson.D{{Key: "_id", Value: userObjectID}}
	update = bson.D{{Key: "$pull", Value: bson.D{{Key: "following", Value: followedUserObjectID}}}}
	if _, err := u.userRepo.UpdateOne(ctx, filter, update); err != nil {
		return nil, constants.ErrInternalServer
	}

	filter = bson.D{{Key: "_id", Value: userObjectID}}
	opts := options.FindOne()
	opts.SetProjection(bson.D{{Key: "hash", Value: 0}})

	user, err := u.userRepo.GetOne(ctx, filter, opts)
	if err != nil {
		return nil, constants.ErrInternalServer
	}

	return user, nil
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
