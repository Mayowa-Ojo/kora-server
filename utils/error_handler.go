package utils

import (
	"github.com/Mayowa-Ojo/kora/types"
	"github.com/gofiber/fiber"
)

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
