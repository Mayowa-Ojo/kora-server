package controllers

import (
	"github.com/Mayowa-Ojo/kora-server/domain"
	"github.com/Mayowa-Ojo/kora-server/utils"
	"github.com/gofiber/fiber"
)

// UserController -
type UserController struct {
	userService domain.UserService
}

// NewUserController -
func NewUserController(u domain.UserService) *UserController {
	return &UserController{
		u,
	}
}

// GetAll - fetch all users from DB collection
func (u *UserController) GetAll(ctx *fiber.Ctx) {
	users, err := u.userService.GetAll(ctx)
	if err != nil {
		ctx.Next(err)

		return
	}

	r := utils.NewResponse()
	r.JSONResponse(ctx, true, fiber.StatusOK, "[INFO]: Resource found", users)
}

// GetOne - fetch user with matching query [e.g id] from DB collection
func (u *UserController) GetOne(ctx *fiber.Ctx) {
	user, err := u.userService.GetOne(ctx)
	if err != nil {
		ctx.Next(err)

		return
	}

	r := utils.NewResponse()
	r.JSONResponse(ctx, true, fiber.StatusOK, "[INFO]: Resource found", user)
}

// GetUserProfile - fetch user profile with matching query [username]
func (u *UserController) GetUserProfile(ctx *fiber.Ctx) {
	user, err := u.userService.GetUserProfile(ctx)
	if err != nil {
		ctx.Next(err)

		return
	}

	r := utils.NewResponse()
	r.JSONResponse(ctx, true, fiber.StatusOK, "[INFO]: Resource found", user)
}

// UpdateProfile -
func (u *UserController) UpdateProfile(ctx *fiber.Ctx) {
	user, err := u.userService.UpdateProfile(ctx)
	if err != nil {
		ctx.Next(err)

		return
	}

	r := utils.NewResponse()
	r.JSONResponse(ctx, true, fiber.StatusOK, "[INFO]: Resource updated", user)
}

// UpdateUserKnowledge -
func (u *UserController) UpdateUserKnowledge(ctx *fiber.Ctx) {
	topics, err := u.userService.UpdateUserKnowledge(ctx)
	if err != nil {
		ctx.Next(err)

		return
	}

	r := utils.NewResponse()
	r.JSONResponse(ctx, true, fiber.StatusOK, "[INFO]: Resource updated", topics)
}

// GetFollowersForUser -
func (u *UserController) GetFollowersForUser(ctx *fiber.Ctx) {
	followers, err := u.userService.GetFollowersForUser(ctx)
	if err != nil {
		ctx.Next(err)

		return
	}

	r := utils.NewResponse()
	r.JSONResponse(ctx, true, fiber.StatusOK, "[INFO]: Resource found", followers)
}

// GetFollowingForUser -
func (u *UserController) GetFollowingForUser(ctx *fiber.Ctx) {
	following, err := u.userService.GetFollowingForUser(ctx)
	if err != nil {
		ctx.Next(err)

		return
	}

	r := utils.NewResponse()
	r.JSONResponse(ctx, true, fiber.StatusOK, "[INFO]: Resource found", following)
}

// GetPostsForUser -
func (u *UserController) GetPostsForUser(ctx *fiber.Ctx) {
	posts, err := u.userService.GetPostsForUser(ctx)
	if err != nil {
		ctx.Next(err)

		return
	}

	r := utils.NewResponse()
	r.JSONResponse(ctx, true, fiber.StatusOK, "[INFO]: Resource found", posts)
}

// GetSharedPostsForUser -
func (u *UserController) GetSharedPostsForUser(ctx *fiber.Ctx) {
	sharedPosts, err := u.userService.GetSharedPostsForUser(ctx)
	if err != nil {
		ctx.Next(err)

		return
	}

	r := utils.NewResponse()
	r.JSONResponse(ctx, true, fiber.StatusOK, "[INFO]: Resource found", sharedPosts)
}

// GetKnowledgeForUser -
func (u *UserController) GetKnowledgeForUser(ctx *fiber.Ctx) {
	knowledge, err := u.userService.GetKnowledgeForUser(ctx)
	if err != nil {
		ctx.Next(err)

		return
	}

	r := utils.NewResponse()
	r.JSONResponse(ctx, true, fiber.StatusOK, "[INFO]: Resource found", knowledge)
}

// GetSpacesForUser -
func (u *UserController) GetSpacesForUser(ctx *fiber.Ctx) {
	spaces, err := u.userService.GetSpacesForUser(ctx)
	if err != nil {
		ctx.Next(err)

		return
	}

	r := utils.NewResponse()
	r.JSONResponse(ctx, true, fiber.StatusOK, "[INFO]: Resource found", spaces)
}

// FollowUser -
func (u *UserController) FollowUser(ctx *fiber.Ctx) {
	user, err := u.userService.FollowUser(ctx)
	if err != nil {
		ctx.Next(err)

		return
	}

	r := utils.NewResponse()
	r.JSONResponse(ctx, true, fiber.StatusOK, "[INFO]: Resource updated", user)
}

// UnfollowUser -
func (u *UserController) UnfollowUser(ctx *fiber.Ctx) {
	user, err := u.userService.UnfollowUser(ctx)
	if err != nil {
		ctx.Next(err)

		return
	}

	r := utils.NewResponse()
	r.JSONResponse(ctx, true, fiber.StatusOK, "[INFO]: Resource updated", user)
}

// SetPinnedPost -
func (u *UserController) SetPinnedPost(ctx *fiber.Ctx) {
	if err := u.userService.SetPinnedPost(ctx); err != nil {
		ctx.Next(err)

		return
	}

	r := utils.NewResponse()
	r.JSONResponse(ctx, true, fiber.StatusOK, "[INFO]: Resource updated", nil)
}

// UnsetPinnedPost -
func (u *UserController) UnsetPinnedPost(ctx *fiber.Ctx) {
	if err := u.userService.UnsetPinnedPost(ctx); err != nil {
		ctx.Next(err)

		return
	}

	r := utils.NewResponse()
	r.JSONResponse(ctx, true, fiber.StatusOK, "[INFO]: Resource updated", nil)
}
