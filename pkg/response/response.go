package response

import (
	"github.com/gofiber/fiber/v2"
	"domesv2/pkg/errors"
)

type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
	Details string      `json:"details,omitempty"`
}

func Success(c *fiber.Ctx, data interface{}, message string) error {
	return c.Status(fiber.StatusOK).JSON(APIResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func Created(c *fiber.Ctx, data interface{}, message string) error {
	return c.Status(fiber.StatusCreated).JSON(APIResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func Error(c *fiber.Ctx, err error) error {
	if appErr, ok := err.(*errors.AppError); ok {
		return c.Status(appErr.Code).JSON(APIResponse{
			Success: false,
			Message: appErr.Message,
			Error:   appErr.Details,
			Details: appErr.Error(),
		})
	}

	return c.Status(fiber.StatusInternalServerError).JSON(APIResponse{
		Success: false,
		Message: "Internal Server Error",
		Error:   "INTERNAL_ERROR",
		Details: err.Error(),
	})
}

func BadRequest(c *fiber.Ctx, message, details string) error {
	err := errors.NewBadRequestError(message, details)
	return Error(c, err)
}

func Unauthorized(c *fiber.Ctx, message, details string) error {
	err := errors.NewUnauthorizedError(message, details)
	return Error(c, err)
}

func Forbidden(c *fiber.Ctx, message, details string) error {
	err := errors.NewForbiddenError(message, details)
	return Error(c, err)
}

func NotFound(c *fiber.Ctx, message, details string) error {
	err := errors.NewNotFoundError(message, details)
	return Error(c, err)
}

func InternalServerError(c *fiber.Ctx, message, details string) error {
	err := errors.NewInternalServerError(message, details)
	return Error(c, err)
}
