package errors

import (
	"fmt"
	"net/http"
)

type AppError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

func (e *AppError) Error() string {
	return fmt.Sprintf("%s: %s", e.Message, e.Details)
}

// Common error constructors
func NewBadRequestError(message, details string) *AppError {
	return &AppError{
		Code:    http.StatusBadRequest,
		Message: message,
		Details: details,
	}
}

func NewUnauthorizedError(message, details string) *AppError {
	return &AppError{
		Code:    http.StatusUnauthorized,
		Message: message,
		Details: details,
	}
}

func NewForbiddenError(message, details string) *AppError {
	return &AppError{
		Code:    http.StatusForbidden,
		Message: message,
		Details: details,
	}
}

func NewNotFoundError(message, details string) *AppError {
	return &AppError{
		Code:    http.StatusNotFound,
		Message: message,
		Details: details,
	}
}

func NewInternalServerError(message, details string) *AppError {
	return &AppError{
		Code:    http.StatusInternalServerError,
		Message: message,
		Details: details,
	}
}

func NewConflictError(message, details string) *AppError {
	return &AppError{
		Code:    http.StatusConflict,
		Message: message,
		Details: details,
	}
}

func NewValidationError(message, details string) *AppError {
	return &AppError{
		Code:    http.StatusUnprocessableEntity,
		Message: message,
		Details: details,
	}
}

// Error codes for easy reference
const (
	ErrCodeInvalidCredentials = "INVALID_CREDENTIALS"
	ErrCodeTokenExpired       = "TOKEN_EXPIRED"
	ErrCodeInvalidToken       = "INVALID_TOKEN"
	ErrCodeUserExists         = "USER_ALREADY_EXISTS"
	ErrCodeUserNotFound       = "USER_NOT_FOUND"
	ErrCodeCaptchaInvalid     = "CAPTCHA_INVALID"
	ErrCodeCaptchaMissing     = "CAPTCHA_MISSING"
	ErrCodeDatabaseError      = "DATABASE_ERROR"
	ErrCodeRedisError         = "REDIS_ERROR"
	ErrCodeInternalServer     = "INTERNAL_ERROR"
	ErrCodeValidationFailed   = "VALIDATION_FAILED"
	ErrCodeInvalidResetToken  = "INVALID_RESET_TOKEN"
)
