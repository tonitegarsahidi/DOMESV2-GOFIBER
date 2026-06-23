package middleware

import (
	goerrors "errors"

	"github.com/gofiber/fiber/v2"
	apperrors "domesv2/pkg/errors"
	"go.uber.org/zap"
)

func GlobalErrorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	message := "Internal Server Error"
	details := "INTERNAL_ERROR"

	var fiberErr *fiber.Error
	if goerrors.As(err, &fiberErr) {
		code = fiberErr.Code
		message = fiberErr.Message
		details = "HTTP_" + fiberErr.Message
	}

	var appErr *apperrors.AppError
	if goerrors.As(err, &appErr) {
		code = appErr.Code
		message = appErr.Message
		details = appErr.Details
	}

	zap.L().Error("Unhandled error",
		zap.Int("status", code),
		zap.String("path", c.Path()),
		zap.String("method", c.Method()),
		zap.String("ip", c.IP()),
		zap.Error(err),
	)

	return c.Status(code).JSON(fiber.Map{
		"success": false,
		"message": message,
		"error":   details,
		"details": err.Error(),
	})
}
