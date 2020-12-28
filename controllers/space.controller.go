package controllers

import (
	"github.com/Mayowa-Ojo/kora/domain"
	"github.com/Mayowa-Ojo/kora/utils"
	"github.com/gofiber/fiber"
)

// SpaceController -
type SpaceController struct {
	spaceService domain.SpaceService
}

// NewSpaceController - presenter layer for spaces
func NewSpaceController(s domain.SpaceService) *SpaceController {
	return &SpaceController{
		s,
	}
}

// GetAll - fetch all spaces from DB collection
func (s *SpaceController) GetAll(ctx *fiber.Ctx) {
	spaces, err := s.spaceService.GetAll(ctx)
	if err != nil {
		ctx.Next(err)

		return
	}

	r := utils.NewResponse()
	r.JSONResponse(ctx, true, fiber.StatusOK, "[INFO]: Resource found", spaces)
}

// GetOne - fetch space with matching query [e.g id] from DB collection
func (s *SpaceController) GetOne(ctx *fiber.Ctx) {
	space, err := s.spaceService.GetOne(ctx)
	if err != nil {
		ctx.Next(err)

		return
	}

	r := utils.NewResponse()
	r.JSONResponse(ctx, true, fiber.StatusOK, "[INFO]: Resource found", space)
}

// GetBySlug - fetch space with matching slug from DB collection
func (s *SpaceController) GetBySlug(ctx *fiber.Ctx) {
	space, err := s.spaceService.GetBySlug(ctx)
	if err != nil {
		ctx.Next(err)

		return
	}

	r := utils.NewResponse()
	r.JSONResponse(ctx, true, fiber.StatusOK, "[INFO]: Resource found", space)
}

// Create - create new space and save to DB collection
func (s *SpaceController) Create(ctx *fiber.Ctx) {
	space, err := s.spaceService.Create(ctx)
	if err != nil {
		ctx.Next(err)

		return
	}

	r := utils.NewResponse()
	r.JSONResponse(ctx, true, fiber.StatusOK, "[INFO]: Resource created", space)
}

// GetPostsForSpace -
func (s *SpaceController) GetPostsForSpace(ctx *fiber.Ctx) {
	posts, err := s.spaceService.GetPostsForSpace(ctx)
	if err != nil {
		ctx.Next(err)

		return
	}

	r := utils.NewResponse()
	r.JSONResponse(ctx, true, fiber.StatusOK, "[INFO]: Resource found", posts)
}

// GetMembersForSpace -
func (s *SpaceController) GetMembersForSpace(ctx *fiber.Ctx) {
	members, err := s.spaceService.GetMembersForSpace(ctx)
	if err != nil {
		ctx.Next(err)

		return
	}

	r := utils.NewResponse()
	r.JSONResponse(ctx, true, fiber.StatusOK, "[INFO]: Resource found", members)
}

// GetSuggestedSpaces -
func (s *SpaceController) GetSuggestedSpaces(ctx *fiber.Ctx) {
	spaces, err := s.spaceService.GetSuggestedSpaces(ctx)
	if err != nil {
		ctx.Next(err)

		return
	}

	r := utils.NewResponse()
	r.JSONResponse(ctx, true, fiber.StatusOK, "[INFO]: Resource found", spaces)
}

// UpdateProfileByAdmin -
func (s *SpaceController) UpdateProfileByAdmin(ctx *fiber.Ctx) {
	space, err := s.spaceService.UpdateProfileByAdmin(ctx)
	if err != nil {
		ctx.Next(err)

		return
	}

	r := utils.NewResponse()
	r.JSONResponse(ctx, true, fiber.StatusOK, "[INFO]: Resource updated", space)
}

// DeleteSpaceByAdmin -
func (s *SpaceController) DeleteSpaceByAdmin(ctx *fiber.Ctx) {
	if err := s.spaceService.DeleteSpaceByAdmin(ctx); err != nil {
		ctx.Next(err)

		return
	}

	r := utils.NewResponse()
	r.JSONResponse(ctx, true, fiber.StatusNoContent, "[INFO]: Resource deleted", nil)
}

// FollowSpace -
func (s *SpaceController) FollowSpace(ctx *fiber.Ctx) {
	if err := s.spaceService.FollowSpace(ctx); err != nil {
		ctx.Next(err)

		return
	}

	r := utils.NewResponse()
	r.JSONResponse(ctx, true, fiber.StatusOK, "[INFO]: Resource updated", nil)
}

// UnfollowSpace -
func (s *SpaceController) UnfollowSpace(ctx *fiber.Ctx) {
	if err := s.spaceService.UnfollowSpace(ctx); err != nil {
		ctx.Next(err)

		return
	}

	r := utils.NewResponse()
	r.JSONResponse(ctx, true, fiber.StatusOK, "[INFO]: Resource updated", nil)
}

// SetPinnedPost -
func (s *SpaceController) SetPinnedPost(ctx *fiber.Ctx) {
	if err := s.spaceService.SetPinnedPost(ctx); err != nil {
		ctx.Next(err)

		return
	}

	r := utils.NewResponse()
	r.JSONResponse(ctx, true, fiber.StatusOK, "[INFO]: Resource updated", nil)
}

// UnsetPinnedPost -
func (s *SpaceController) UnsetPinnedPost(ctx *fiber.Ctx) {
	if err := s.spaceService.UnsetPinnedPost(ctx); err != nil {
		ctx.Next(err)

		return
	}

	r := utils.NewResponse()
	r.JSONResponse(ctx, true, fiber.StatusOK, "[INFO]: Resource updated", nil)
}
