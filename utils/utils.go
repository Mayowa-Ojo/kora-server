package utils

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/Mayowa-Ojo/kora/constants"
	"github.com/Mayowa-Ojo/kora/domain"
	"github.com/Mayowa-Ojo/kora/entity"
	"github.com/Mayowa-Ojo/kora/types"
	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber"
	sid "github.com/ventu-io/go-shortid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

// Response -
type Response struct {
	Ok      bool      `json:"ok"`
	Message string    `json:"message"`
	Code    int       `json:"code"`
	Data    types.Any `json:"data"`
}

// NewResponse -
func NewResponse() *Response {
	return &Response{}
}

// JSONResponse -
func (r Response) JSONResponse(c *fiber.Ctx, ok bool, code int, message string, data types.Any) {
	r.Ok = ok
	r.Code = code
	r.Message = message
	r.Data = data

	if err := c.Status(code).JSON(r); err != nil {
		c.Next(err)
		return
	}
}

// GenerateID - generate a short unique id
func GenerateID() (string, error) {
	return sid.Generate()
}

// ValidatePassword -
func ValidatePassword(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))

	return err == nil
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

// GenerateUsername - creates a unique username from the user's first and last names
func GenerateUsername(names map[string]string, ctx *fiber.Ctx, userRepo domain.UserRepository) (string, error) {
	var username string

	if _, ok := names["firstname"]; !ok {
		return "", errors.New("missing firstname field")
	}

	if _, ok := names["lastname"]; !ok {
		return "", errors.New("missing lastname field")
	}

	filter := bson.D{{
		Key: "$and",
		Value: bson.A{
			bson.D{{Key: "firstname", Value: names["firstname"]}},
			bson.D{{Key: "lastname", Value: names["lastname"]}},
		},
	}}
	opts := options.Find()
	users, err := userRepo.GetMany(ctx, filter, opts)

	if err != nil {
		return "", constants.ErrInternalServer
	}

	matches := len(users)

	if matches < 1 {
		username = names["firstname"] + "-" + names["lastname"]

		return username, nil
	}

	// username = names["firstname"] + "-" + names["lastname"] + "-" + strconv.Itoa(matches)
	username = fmt.Sprintf("%s-%s-%s", names["firstname"], names["lastname"], strconv.Itoa(matches))
	return username, nil
}

// ErrorHandler -
func ErrorHandler(ctx *fiber.Ctx, err error) {
	code := fiber.StatusInternalServerError
	r := NewResponse()

	if e, ok := err.(*fiber.Error); ok {
		r.JSONResponse(ctx, false, e.Code, e.Message, make(types.GenericMap))
	} else {
		r.JSONResponse(ctx, false, code, "[Error]: Internal server error", make(types.GenericMap))
	}
}

// JWTError -
func JWTError(ctx *fiber.Ctx, err error) {
	if err.Error() == "Missing or malformed JWT" {
		ctx.Next(constants.ErrUnprocessableEntity)
	}

	ctx.Next(constants.ErrInvalidCredentials)
}

// GenerateSlug - creates a url friendly slug from a string s
//               [e.g] 'What are generics?' becomes 'what-are-generics'
func GenerateSlug(s string) string {
	s = strings.ToLower(s)
	slice := strings.Split(s, " ")
	replacer := strings.NewReplacer(
		".", "", ",", "", "&", "", "@", "", "#", "", "$", "", "*", "", "+", "",
	)

	s = strings.Join(slice, "-")
	s = strings.TrimSuffix(s, "?")
	s = replacer.Replace(s)

	return s
}
