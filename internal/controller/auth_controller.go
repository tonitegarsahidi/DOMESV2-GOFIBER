package controller

import (
	"strings"

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

func (ctrl *AuthController) UpdateProfile(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)

	var req model.UpdateProfileRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request body", "INVALID_REQUEST_BODY")
	}

	result, err := ctrl.authService.UpdateProfile(userID, &req)
	if err != nil {
		return response.Error(c, err)
	}

	return response.Success(c, result, "Profile updated successfully")
}

func (ctrl *AuthController) ChangePassword(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)

	var req model.ChangePasswordRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request body", "INVALID_REQUEST_BODY")
	}

	if err := ctrl.authService.ChangePassword(userID, &req); err != nil {
		return response.Error(c, err)
	}

	return response.Success(c, nil, "Password changed successfully")
}

func (ctrl *AuthController) GetNotificationPreferences(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)

	result, err := ctrl.authService.GetNotificationPreferences(userID)
	if err != nil {
		return response.Error(c, err)
	}

	return response.Success(c, result, "Notification preferences retrieved successfully")
}

func (ctrl *AuthController) UpdateNotificationPreferences(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)

	var req model.UpdateNotificationRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request body", "INVALID_REQUEST_BODY")
	}

	result, err := ctrl.authService.UpdateNotificationPreferences(userID, &req)
	if err != nil {
		return response.Error(c, err)
	}

	return response.Success(c, result, "Notification preferences updated successfully")
}

func (ctrl *AuthController) GetAdminEmails(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)
	profile, err := ctrl.authService.GetProfile(userID)
	if err != nil {
		return response.Error(c, err)
	}
	if profile.Role == nil || strings.ToLower(*profile.Role) != "administrator" {
		return response.Forbidden(c, "Only administrators can access this resource", "FORBIDDEN")
	}

	result, err := ctrl.authService.GetAdminEmails()
	if err != nil {
		return response.Error(c, err)
	}

	return response.Success(c, result, "Admin emails retrieved successfully")
}

func (ctrl *AuthController) AddAdminEmail(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)
	profile, err := ctrl.authService.GetProfile(userID)
	if err != nil {
		return response.Error(c, err)
	}
	if profile.Role == nil || strings.ToLower(*profile.Role) != "administrator" {
		return response.Forbidden(c, "Only administrators can perform this action", "FORBIDDEN")
	}

	var req model.AddAdminEmailRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request body", "INVALID_REQUEST_BODY")
	}

	result, err := ctrl.authService.AddAdminEmail(req.Email)
	if err != nil {
		return response.Error(c, err)
	}

	return response.Created(c, result, "Admin email added successfully")
}

func (ctrl *AuthController) DeleteAdminEmail(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)
	profile, err := ctrl.authService.GetProfile(userID)
	if err != nil {
		return response.Error(c, err)
	}
	if profile.Role == nil || strings.ToLower(*profile.Role) != "administrator" {
		return response.Forbidden(c, "Only administrators can perform this action", "FORBIDDEN")
	}

	email := c.Params("email")
	if email == "" {
		return response.BadRequest(c, "Email parameter is required", "VALIDATION_FAILED")
	}

	if err := ctrl.authService.DeleteAdminEmail(email); err != nil {
		return response.Error(c, err)
	}

	return response.Success(c, nil, "Admin email removed successfully")
}
