package controllers

import (
	"github.com/Mayowa-Ojo/kora/domain"
	"github.com/Mayowa-Ojo/kora/utils"
	"github.com/gofiber/fiber"
)

// TopicController -
type TopicController struct {
	topicService domain.TopicService
}

// NewTopicController -
func NewTopicController(t domain.TopicService) *TopicController {
	return &TopicController{
		t,
	}
}

// GetAll - fetch all topics from DB collection
func (t *TopicController) GetAll(ctx *fiber.Ctx) {
	topics, err := t.topicService.GetAll(ctx)
	if err != nil {
		ctx.Next(err)

		return
	}

	r := utils.NewResponse()
	r.JSONResponse(ctx, true, fiber.StatusFound, "[INFO]: Resource found", topics)
}

// GetOne - fetch topic with matching query [e.g id] from DB collection
func (t *TopicController) GetOne(ctx *fiber.Ctx) {
	topic, err := t.topicService.GetOne(ctx)
	if err != nil {
		ctx.Next(err)

		return
	}

	r := utils.NewResponse()
	r.JSONResponse(ctx, true, fiber.StatusFound, "[INFO]: Resource found", topic)
}

// Create - create new topic and save to DB collection
func (t *TopicController) Create(ctx *fiber.Ctx) {
	topic, err := t.topicService.Create(ctx)
	if err != nil {
		ctx.Next(err)

		return
	}

	r := utils.NewResponse()
	r.JSONResponse(ctx, true, fiber.StatusOK, "[INFO]: Resource created", topic)
}
