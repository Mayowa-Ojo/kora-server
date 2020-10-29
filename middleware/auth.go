package middleware

import (
	"github.com/Mayowa-Ojo/kora/config"
	"github.com/Mayowa-Ojo/kora/constants"
	"github.com/Mayowa-Ojo/kora/entity"
	"github.com/Mayowa-Ojo/kora/utils"
	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber"
	jwtware "github.com/gofiber/jwt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// AuthorizeRoute -
func AuthorizeRoute() fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningKey:   []byte("jwt-secret"),
		ErrorHandler: utils.JWTError,
	})
}

// ExtractUserFromAuthHeader -
func ExtractUserFromAuthHeader(conn *config.DBConn) fiber.Handler {
	return func(ctx *fiber.Ctx) {
		user := ctx.Locals("user").(*jwt.Token)
		claims := user.Claims.(jwt.MapClaims)

		if _, ok := claims["userId"].(string); !ok {
			ctx.Next(constants.ErrUnauthorized)
			return
		}

		userID := claims["userId"].(string)
		userObjectID, err := primitive.ObjectIDFromHex(userID)
		if err != nil {
			ctx.Next(constants.ErrInternalServer)
			return
		}

		authUser := new(entity.User)
		filter := bson.D{{Key: "_id", Value: userObjectID}}
		opts := options.FindOne()
		result := conn.DB.Collection("users").FindOne(ctx.Fasthttp, filter, opts)

		if err := result.Decode(&authUser); err != nil {
			ctx.Next(constants.ErrInternalServer)
			return
		}

		ctx.Locals("authUser", authUser)
		ctx.Next()
	}
}
