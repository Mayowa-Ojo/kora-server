package constants

import (
	"github.com/gofiber/fiber"
)

var (
	// ErrResourceExists -
	ErrResourceExists = &fiber.Error{
		Code:    fiber.ErrBadRequest.Code,
		Message: "resource already exists",
	}

	// ErrNotFound -
	ErrNotFound = &fiber.Error{
		Code:    fiber.ErrNotFound.Code,
		Message: "resource doesn't exist",
	}

	// ErrUnprocessableEntity -
	ErrUnprocessableEntity = &fiber.Error{
		Code:    fiber.ErrUnprocessableEntity.Code,
		Message: "error parsing request body",
	}

	// ErrInvalidCredentials -
	ErrInvalidCredentials = &fiber.Error{
		Code:    fiber.ErrPreconditionFailed.Code,
		Message: "invalid credentials",
	}

	// ErrUnauthorized -
	ErrUnauthorized = &fiber.Error{
		Code:    fiber.ErrUnauthorized.Code,
		Message: "client is not authenticated",
	}

	// ErrForbidden -
	ErrForbidden = &fiber.Error{
		Code:    fiber.ErrForbidden.Code,
		Message: "client is restricted from viewing this resource",
	}

	// ErrInternalServer -
	ErrInternalServer = &fiber.Error{
		Code:    fiber.ErrInternalServerError.Code,
		Message: "internal server error",
	}
)
