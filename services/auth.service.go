package services

import (
	"time"

	"github.com/Mayowa-Ojo/kora/config"
	"github.com/Mayowa-Ojo/kora/constants"

	"github.com/Mayowa-Ojo/kora/domain"
	"github.com/Mayowa-Ojo/kora/entity"
	"github.com/Mayowa-Ojo/kora/types"
	"github.com/Mayowa-Ojo/kora/utils"
	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber"
	"go.mongodb.org/mongo-driver/bson"
	mg "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	env = config.NewEnvConfig()
)

// Auth -
type Auth struct {
	userRepo domain.UserRepository
}

// NewAuthService -
func NewAuthService(a domain.UserRepository) domain.AuthService {
	return &Auth{
		a,
	}
}

// Login -
func (a Auth) Login(ctx *fiber.Ctx) (types.GenericMap, error) {
	var credentials struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := ctx.BodyParser(&credentials); err != nil {
		return nil, constants.ErrUnprocessableEntity
	}

	filter := bson.D{{Key: "email", Value: credentials.Email}}
	opts := options.FindOne()
	opts.SetProjection(bson.D{{}})
	user, err := a.userRepo.GetOne(ctx, filter, opts)
	if err != nil {
		if err == mg.ErrNoDocuments {
			return nil, constants.ErrNotFound
		}
		return nil, constants.ErrInternalServer
	}

	if isValid := utils.ValidatePassword(credentials.Password, user.Hash); !isValid {
		return nil, constants.ErrInvalidCredentials
	}

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["userId"] = user.ID
	claims["username"] = user.Username
	claims["exp"] = time.Now().Add(time.Hour * 168).Unix()

	t, err := token.SignedString([]byte(env.JwtSecret))
	if err != nil {
		return nil, constants.ErrInternalServer
	}

	result := types.GenericMap{
		"user":  user,
		"token": t,
	}

	return result, nil
}

// Signup -
func (a Auth) Signup(ctx *fiber.Ctx) (types.GenericMap, error) {
	var credentials struct {
		Firstname string `json:"firstname"`
		Lastname  string `json:"lastname"`
		Email     string `json:"email"`
		Password  string `json:"password"`
	}

	if err := ctx.BodyParser(&credentials); err != nil {
		return nil, constants.ErrUnprocessableEntity
	}

	filter := bson.D{{Key: "email", Value: credentials.Email}}
	opts := options.FindOne()
	user, err := a.userRepo.GetOne(ctx, filter, opts)
	if err != nil && err != mg.ErrNoDocuments {
		return nil, constants.ErrInternalServer
	}
	if user != nil {
		return nil, constants.ErrResourceExists
	}

	username, err := utils.GenerateUsername(
		map[string]string{
			"firstname": credentials.Firstname,
			"lastname":  credentials.Lastname,
		},
		ctx,
		a.userRepo,
	)

	if err != nil {
		return nil, constants.ErrInternalServer
	}

	hash, err := utils.GeneratePasswordHash(credentials.Password)
	if err != nil {
		return nil, constants.ErrInternalServer
	}

	instance := entity.User{
		Firstname: credentials.Firstname,
		Lastname:  credentials.Lastname,
		Email:     credentials.Email,
		Username:  username,
		Hash:      hash,
	}

	if err := instance.Validate(); err != nil {
		return nil, constants.ErrInvalidCredentials
	}

	instance.SetDefaultValues()

	insertResult, err := a.userRepo.Create(ctx, instance)
	if err != nil {
		return nil, constants.ErrInternalServer
	}

	filter = bson.D{{Key: "_id", Value: insertResult.InsertedID}}
	opts = options.FindOne()
	user, err = a.userRepo.GetOne(ctx, filter, opts)
	if err != nil {
		return nil, constants.ErrInternalServer
	}

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["userId"] = user.ID
	claims["username"] = user.Username
	claims["exp"] = time.Now().Add(time.Hour * 168).Unix()

	t, err := token.SignedString([]byte(env.JwtSecret))
	if err != nil {
		return nil, constants.ErrInternalServer
	}

	result := types.GenericMap{
		"user":  user,
		"token": t,
	}

	return result, nil
}
