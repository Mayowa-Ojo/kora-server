package controllers

import (
	"github.com/Mayowa-Ojo/kora-server/domain"
	"github.com/Mayowa-Ojo/kora-server/utils"
	"github.com/gofiber/fiber"
)

// ImageController -
type ImageController struct {
	service domain.ImageService
}

// NewImageController -
func NewImageController(i domain.ImageService) *ImageController {
	return &ImageController{
		i,
	}
}

// UploadImage -
func (i *ImageController) UploadImage(ctx *fiber.Ctx) {
	result, err := i.service.UploadImage(ctx)
	if err != nil {
		ctx.Next(err)

		return
	}

	r := utils.NewResponse()
	r.JSONResponse(ctx, true, fiber.StatusCreated, "[INFO]: Resource created", result)
}
