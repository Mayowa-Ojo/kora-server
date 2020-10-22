package utils

import (
	"errors"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber"
	sid "github.com/ventu-io/go-shortid"
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
