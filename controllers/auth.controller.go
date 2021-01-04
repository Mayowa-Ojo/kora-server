package controllers

import (
	"github.com/Mayowa-Ojo/kora-server/domain"
	"github.com/Mayowa-Ojo/kora-server/utils"
	"github.com/gofiber/fiber"
)

// Auth -
type Auth struct {
	service domain.AuthService
}

// NewAuthController -
func NewAuthController(s domain.AuthService) *Auth {
	return &Auth{
		s,
	}
}

// Login -
func (a Auth) Login(ctx *fiber.Ctx) {
	user, err := a.service.Login(ctx)
	if err != nil {
		ctx.Next(err)

		return
	}

	r := utils.NewResponse()
	r.JSONResponse(ctx, true, fiber.StatusOK, "[INFO]: Authentication successful", user)
}

// Signup -
func (a Auth) Signup(ctx *fiber.Ctx) {
	user, err := a.service.Signup(ctx)
	if err != nil {
		ctx.Next(err)

		return
	}

	r := utils.NewResponse()
	r.JSONResponse(ctx, true, fiber.StatusOK, "[INFO]: Signup successful", user)
}
