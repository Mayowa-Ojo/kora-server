package middleware

import (
	"github.com/Mayowa-Ojo/kora/utils"
	"github.com/gofiber/fiber"
	jwtware "github.com/gofiber/jwt"
)

// AuthorizeRoute -
func AuthorizeRoute() fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningKey:   []byte("jwt-secret"),
		ErrorHandler: utils.JWTError,
	})
}
