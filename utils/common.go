package utils

import (
	"errors"

	"github.com/Mayowa-Ojo/kora/constants"
	"github.com/Mayowa-Ojo/kora/domain"
	"github.com/Mayowa-Ojo/kora/entity"
	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber"
	sid "github.com/ventu-io/go-shortid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

// GenerateID - generate a short unique id
func GenerateID() (string, error) {
	return sid.Generate()
}

// GeneratePasswordHash - create secure hash for password input
func GeneratePasswordHash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 12)

	if err != nil {
		return "", err
	}

	return string(hash), nil
}

// GetJwtClaims - extract user details from claims
func GetJwtClaims(ctx *fiber.Ctx, key string) (string, error) {
	user := ctx.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)

	if _, ok := claims[key].(string); !ok {
		return "", errors.New("invalid key")
	}

	return claims[key].(string), nil
}

// GetUserFromAuthHeader -
func GetUserFromAuthHeader(ctx *fiber.Ctx, userRepo domain.UserRepository) (*entity.User, error) {
	userID, err := GetJwtClaims(ctx, "userId")
	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, constants.ErrUnauthorized
	}

	filter := bson.D{{Key: "_id", Value: userObjectID}}
	opts := options.FindOne()
	user, err := userRepo.GetOne(ctx, filter, opts)
	if err != nil {
		return nil, constants.ErrInternalServer
	}

	return user, nil
}
