package controllers

import (
	"github.com/Mayowa-Ojo/kora/domain"
	"github.com/Mayowa-Ojo/kora/utils"
	"github.com/gofiber/fiber"
)

// CommentController -
type CommentController struct {
	commentService domain.CommentService
}

// NewCommentController -
func NewCommentController(c domain.CommentService) *CommentController {
	return &CommentController{
		c,
	}
}

// GetAll - fetch all comments from DB collection
func (c *CommentController) GetAll(ctx *fiber.Ctx) {
	comments, err := c.commentService.GetAll(ctx)
	if err != nil {
		ctx.Next(err)

		return
	}

	r := utils.NewResponse()
	r.JSONResponse(ctx, true, fiber.StatusFound, "[INFO]: Resource found", comments)
}

// GetOne - fetch comment with matching query [e.g id] from DB collection
func (c *CommentController) GetOne(ctx *fiber.Ctx) {
	comment, err := c.commentService.GetOne(ctx)
	if err != nil {
		ctx.Next(err)

		return
	}

	r := utils.NewResponse()
	r.JSONResponse(ctx, true, fiber.StatusFound, "[INFO]: Resource found", comment)
}

// Create - create new comment and save to DB collection
func (c *CommentController) Create(ctx *fiber.Ctx) {
	comment, err := c.commentService.Create(ctx)
	if err != nil {
		ctx.Next(err)

		return
	}

	r := utils.NewResponse()
	r.JSONResponse(ctx, true, fiber.StatusOK, "[INFO]: Resource created", comment)
}
