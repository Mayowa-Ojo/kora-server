package services

import (
	"time"

	"github.com/Mayowa-Ojo/kora/constants"

	"github.com/Mayowa-Ojo/kora/domain"
	"github.com/Mayowa-Ojo/kora/entity"
	"github.com/Mayowa-Ojo/kora/types"
	"github.com/Mayowa-Ojo/kora/utils"
	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber"
	mg "go.mongodb.org/mongo-driver/mongo"
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

	user, err := a.userRepo.GetOne(ctx, types.GenericMap{"email": credentials.Email})
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

	t, err := token.SignedString([]byte("jwt-secret"))
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
		Username  string `json:"username"`
		Email     string `json:"email"`
		Password  string `json:"password"`
	}

	if err := ctx.BodyParser(&credentials); err != nil {
		return nil, constants.ErrUnprocessableEntity
	}

	user, err := a.userRepo.GetOne(ctx, types.GenericMap{"email": credentials.Email})
	if err != nil && err != mg.ErrNoDocuments {
		return nil, constants.ErrInternalServer
	}
	if user != nil {
		return nil, constants.ErrResourceExists
	}

	hash, err := utils.GeneratePasswordHash(credentials.Password)
	if err != nil {
		return nil, constants.ErrInternalServer
	}

	instance := entity.User{
		Firstname: credentials.Firstname,
		Lastname:  credentials.Lastname,
		Email:     credentials.Email,
		Username:  credentials.Username,
		Hash:      hash,
	}

	if err := instance.Validate(); err != nil {
		return nil, constants.ErrInvalidCredentials
	}

	insertResult, err := a.userRepo.Create(ctx, instance)
	if err != nil {
		return nil, constants.ErrInternalServer
	}

	user, err = a.userRepo.GetOne(ctx, types.GenericMap{"_id": insertResult.InsertedID})
	if err != nil {
		return nil, constants.ErrInternalServer
	}

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["userId"] = user.ID
	claims["username"] = user.Username
	claims["exp"] = time.Now().Add(time.Hour * 168).Unix()

	t, err := token.SignedString([]byte("jwt secret"))
	if err != nil {
		return nil, constants.ErrInternalServer
	}

	result := types.GenericMap{
		"user":  user,
		"token": t,
	}

	return result, nil
}