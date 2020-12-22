package controllers

import (
	"github.com/Mayowa-Ojo/kora/domain"
	"github.com/Mayowa-Ojo/kora/utils"
	"github.com/gofiber/fiber"
)

// PostController - Structure for a post controller
type PostController struct {
	postService domain.PostService
	userService domain.UserService
}

// NewPostController - Creates post controller instance
func NewPostController(p domain.PostService, u domain.UserService) *PostController {
	return &PostController{
		p,
		u,
	}
}

// GetAll - fetch all posts from DB collection
func (p *PostController) GetAll(ctx *fiber.Ctx) {
	posts, err := p.postService.GetAll(ctx)
	if err != nil {
		ctx.Next(err)

		return
	}

	r := utils.NewResponse()
	r.JSONResponse(ctx, true, fiber.StatusOK, "[INFO]: Resource found", posts)
}

// GetOne - fetch post with matching query [e.g id] from DB collection
func (p *PostController) GetOne(ctx *fiber.Ctx) {
	post, err := p.postService.GetOne(ctx)
	if err != nil {
		ctx.Next(err)

		return
	}

	err = p.userService.UpdateContentViews(ctx)
	if err != nil {
		ctx.Next(err)

		return
	}

	r := utils.NewResponse()
	r.JSONResponse(ctx, true, fiber.StatusOK, "[INFO]: Resource found", post)
}

// GetBySlug - fetch post with matching query [e.g id] from DB collection
func (p *PostController) GetBySlug(ctx *fiber.Ctx) {
	post, err := p.postService.GetBySlug(ctx)
	if err != nil {
		ctx.Next(err)

		return
	}

	err = p.userService.UpdateContentViews(ctx)
	if err != nil {
		ctx.Next(err)

		return
	}

	r := utils.NewResponse()
	r.JSONResponse(ctx, true, fiber.StatusOK, "[INFO]: Resource found", post)
}

// Create - create new post and save to DB collection
func (p *PostController) Create(ctx *fiber.Ctx) {
	post, err := p.postService.Create(ctx)
	if err != nil {
		ctx.Next(err)

		return
	}

	r := utils.NewResponse()
	r.JSONResponse(ctx, true, fiber.StatusCreated, "[INFO]: Resource created", post)
}

// DeleteOne - create new post and save to DB collection
func (p *PostController) DeleteOne(ctx *fiber.Ctx) {
	if err := p.postService.DeleteOne(ctx); err != nil {
		ctx.Next(err)

		return
	}

	r := utils.NewResponse()
	r.JSONResponse(ctx, true, fiber.StatusOK, "[INFO]: Resource deleted", nil)
}

// GetFeedForUser -
func (p *PostController) GetFeedForUser(ctx *fiber.Ctx) {
	posts, err := p.postService.GetFeedForUser(ctx)
	if err != nil {
		ctx.Next(err)

		return
	}

	r := utils.NewResponse()
	r.JSONResponse(ctx, true, fiber.StatusOK, "[INFO]: Resource found", posts)
}

// GetQuestionsForUser -
func (p *PostController) GetQuestionsForUser(ctx *fiber.Ctx) {
	posts, err := p.postService.GetQuestionsForUser(ctx)
	if err != nil {
		ctx.Next(err)

		return
	}

	r := utils.NewResponse()
	r.JSONResponse(ctx, true, fiber.StatusOK, "[INFO]: Resource found", posts)
}

// UpvotePostByUser -
func (p *PostController) UpvotePostByUser(ctx *fiber.Ctx) {
	if err := p.postService.UpvotePostByUser(ctx); err != nil {
		ctx.Next(err)

		return
	}

	r := utils.NewResponse()
	r.JSONResponse(ctx, true, fiber.StatusOK, "[INFO]: Resource updated", nil)
}

// DownvotePostByUser -
func (p *PostController) DownvotePostByUser(ctx *fiber.Ctx) {
	if err := p.postService.DownvotePostByUser(ctx); err != nil {
		ctx.Next(err)

		return
	}

	r := utils.NewResponse()
	r.JSONResponse(ctx, true, fiber.StatusOK, "[INFO]: Resource updated", nil)
}

// FollowPost -
func (p *PostController) FollowPost(ctx *fiber.Ctx) {
	if err := p.postService.FollowPost(ctx); err != nil {
		ctx.Next(err)

		return
	}

	r := utils.NewResponse()
	r.JSONResponse(ctx, true, fiber.StatusOK, "[INFO]: Resource updated", nil)
}

// UnfollowPost -
func (p *PostController) UnfollowPost(ctx *fiber.Ctx) {
	if err := p.postService.UnfollowPost(ctx); err != nil {
		ctx.Next(err)

		return
	}

	r := utils.NewResponse()
	r.JSONResponse(ctx, true, fiber.StatusOK, "[INFO]: Resource updated", nil)
}

// SharePost - create a shared post
func (p *PostController) SharePost(ctx *fiber.Ctx) {
	sharedPost, err := p.postService.CreateSharedPost(ctx)
	if err != nil {
		ctx.Next(err)

		return
	}

	r := utils.NewResponse()
	r.JSONResponse(ctx, true, fiber.StatusCreated, "[INFO]: Resource created", sharedPost)
}
