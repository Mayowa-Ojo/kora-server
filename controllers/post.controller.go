package controllers

import (
	"strconv"

	"github.com/Mayowa-Ojo/kora/domain"
	"github.com/Mayowa-Ojo/kora/types"
	"github.com/Mayowa-Ojo/kora/utils"
	"github.com/gofiber/fiber"
)

// Post - Structre for a post controller
type Post struct {
	service domain.PostService
}

// NewPostController - Creates post controller instance
func NewPostController(s domain.PostService) *Post {
	return &Post{
		s,
	}
}

// GetAll - Delivers all posts data to the client
func (p *Post) GetAll(ctx *fiber.Ctx) {
	limit := ctx.Query("limit", "10")
	limitInt, err := strconv.ParseInt(limit, 10, 64)

	if err != nil {
		err := new(fiber.Error)
		err.Code = fiber.StatusInternalServerError
		err.Message = "[Error]: Something went wrong"
		ctx.Next(err)

		return
	}

	posts, err := p.service.GetAll(ctx, types.GenericMap{"limit": limitInt})

	if err != nil {
		err := new(fiber.Error)
		err.Code = 404
		err.Message = "[Error]: Resource not found"
		ctx.Next(err)

		return
	}

	r := utils.NewResponse()
	r.JSONResponse(ctx, true, fiber.StatusFound, "[INFO]: Resource found", posts)
}
