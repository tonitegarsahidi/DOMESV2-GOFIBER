package controller

import (
	"github.com/gofiber/fiber/v2"
	"domesv2/internal/model"
	"domesv2/internal/service"
	"domesv2/pkg/response"
)

type AuthController struct {
	authService service.AuthService
}

func NewAuthController(authService service.AuthService) *AuthController {
	return &AuthController{
		authService: authService,
	}
}

func (ctrl *AuthController) Register(c *fiber.Ctx) error {
	var req model.RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request body", "INVALID_REQUEST_BODY")
	}

	result, err := ctrl.authService.Register(&req)
	if err != nil {
		return response.Error(c, err)
	}

	return response.Created(c, result, "User registered successfully")
}

func (ctrl *AuthController) Login(c *fiber.Ctx) error {
	var req model.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request body", "INVALID_REQUEST_BODY")
	}

	result, err := ctrl.authService.Login(&req)
	if err != nil {
		return response.Error(c, err)
	}

	return response.Success(c, result, "Login successful")
}

func (ctrl *AuthController) ForgotPassword(c *fiber.Ctx) error {
	var req model.ForgotPasswordRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request body", "INVALID_REQUEST_BODY")
	}

	if err := ctrl.authService.ForgotPassword(&req); err != nil {
		return response.Error(c, err)
	}

	return response.Success(c, nil, "If the email exists, a reset link has been sent")
}

func (ctrl *AuthController) ResetPassword(c *fiber.Ctx) error {
	var req model.ResetPasswordRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request body", "INVALID_REQUEST_BODY")
	}

	if err := ctrl.authService.ResetPassword(&req); err != nil {
		return response.Error(c, err)
	}

	return response.Success(c, nil, "Password has been reset successfully")
}

func (ctrl *AuthController) Me(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)

	profile, err := ctrl.authService.GetProfile(userID)
	if err != nil {
		return response.Error(c, err)
	}

	return response.Success(c, profile, "User profile retrieved successfully")
}
