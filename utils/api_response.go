package utils

import (
	"github.com/Mayowa-Ojo/kora/types"
	"github.com/gofiber/fiber"
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
