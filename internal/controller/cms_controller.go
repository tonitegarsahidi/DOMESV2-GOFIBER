package controller

import (
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"domesv2/internal/model"
	"domesv2/internal/service"
	"domesv2/pkg/response"
)

type CmsController struct {
	cmsService  service.CmsService
	authService service.AuthService
}

func NewCmsController(cmsService service.CmsService, authService service.AuthService) *CmsController {
	return &CmsController{
		cmsService:  cmsService,
		authService: authService,
	}
}

func (ctrl *CmsController) GetDashboardStats(c *fiber.Ctx) error {
	result, err := ctrl.cmsService.GetDashboardStats()
	if err != nil {
		return response.Error(c, err)
	}
	return response.Success(c, result["stats"], "Dashboard stats retrieved successfully")
}

func (ctrl *CmsController) GetRecentActivity(c *fiber.Ctx) error {
	limit := c.QueryInt("limit", 10)
	result, err := ctrl.cmsService.GetRecentActivity(limit)
	if err != nil {
		return response.Error(c, err)
	}
	return response.Success(c, result, "Recent activity retrieved successfully")
}

func (ctrl *CmsController) GetAnalyticsSummary(c *fiber.Ctx) error {
	period := c.Query("period", "30d")
	result, err := ctrl.cmsService.GetAnalyticsSummary(period)
	if err != nil {
		return response.Error(c, err)
	}
	return response.Success(c, result, "Analytics summary retrieved successfully")
}

func (ctrl *CmsController) GetTopDownloads(c *fiber.Ctx) error {
	limit := c.QueryInt("limit", 10)
	result, err := ctrl.cmsService.GetTopDownloads(limit)
	if err != nil {
		return response.Error(c, err)
	}
	return response.Success(c, result, "Top downloads retrieved successfully")
}

func (ctrl *CmsController) GetTopViews(c *fiber.Ctx) error {
	limit := c.QueryInt("limit", 10)
	result, err := ctrl.cmsService.GetTopViews(limit)
	if err != nil {
		return response.Error(c, err)
	}
	return response.Success(c, result, "Top views retrieved successfully")
}

// User Management (Admin only)
func (ctrl *CmsController) ListUsers(c *fiber.Ctx) error {
	if err := ctrl.verifyAdmin(c); err != nil {
		return response.Error(c, err)
	}

	search := c.Query("search")
	role := c.Query("role")
	status := c.Query("status")
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 20)

	result, err := ctrl.cmsService.ListUsers(search, role, status, page, limit)
	if err != nil {
		return response.Error(c, err)
	}

	return response.Success(c, result, "Users retrieved successfully")
}

func (ctrl *CmsController) CreateUser(c *fiber.Ctx) error {
	if err := ctrl.verifyAdmin(c); err != nil {
		return response.Error(c, err)
	}

	var req model.CreateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request body", "INVALID_REQUEST_BODY")
	}

	user, err := ctrl.cmsService.CreateUser(&req)
	if err != nil {
		return response.Error(c, err)
	}

	return response.Created(c, fiber.Map{
		"id":           user.ID,
		"first_name":   user.FirstName,
		"last_name":    user.LastName,
		"email":        user.Email,
		"organization": user.Organization,
		"position":     user.Position,
		"role":         user.Role,
		"status":       user.Status,
		"created_at":   user.CreatedAt,
	}, "User created successfully")
}

func (ctrl *CmsController) UpdateUser(c *fiber.Ctx) error {
	if err := ctrl.verifyAdmin(c); err != nil {
		return response.Error(c, err)
	}

	idParam := c.Params("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return response.BadRequest(c, "Invalid user ID", "VALIDATION_FAILED")
	}

	var req model.UpdateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request body", "INVALID_REQUEST_BODY")
	}

	user, err := ctrl.cmsService.UpdateUser(uint(id), &req)
	if err != nil {
		return response.Error(c, err)
	}

	return response.Success(c, fiber.Map{
		"id":           user.ID,
		"first_name":   user.FirstName,
		"last_name":    user.LastName,
		"email":        user.Email,
		"organization": user.Organization,
		"position":     user.Position,
		"role":         user.Role,
		"status":       user.Status,
		"updated_at":   user.UpdatedAt,
	}, "User updated successfully")
}

func (ctrl *CmsController) DeleteUser(c *fiber.Ctx) error {
	if err := ctrl.verifyAdmin(c); err != nil {
		return response.Error(c, err)
	}

	idParam := c.Params("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return response.BadRequest(c, "Invalid user ID", "VALIDATION_FAILED")
	}

	if err := ctrl.cmsService.DeleteUser(uint(id)); err != nil {
		return response.Error(c, err)
	}

	return response.Success(c, nil, "User deleted successfully")
}

// Helper to check for admin role
func (ctrl *CmsController) verifyAdmin(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)
	profile, err := ctrl.authService.GetProfile(userID)
	if err != nil {
		return err
	}
	if profile.Role == nil || strings.ToLower(*profile.Role) != "administrator" {
		return response.Forbidden(c, "Only administrators can access this resource", "FORBIDDEN")
	}
	return nil
}

func (ctrl *CmsController) ListReferences(c *fiber.Ctx) error {
	refType := c.Params("type")
	result, err := ctrl.cmsService.ListReferences(refType)
	if err != nil {
		return response.Error(c, err)
	}
	return response.Success(c, result, "Reference list retrieved successfully")
}

func (ctrl *CmsController) CreateReference(c *fiber.Ctx) error {
	if err := ctrl.verifyAdmin(c); err != nil {
		return response.Error(c, err)
	}

	refType := c.Params("type")
	var req model.ReferenceRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request body", "INVALID_REQUEST_BODY")
	}

	result, err := ctrl.cmsService.CreateReference(refType, &req)
	if err != nil {
		return response.Error(c, err)
	}
	return response.Created(c, result, "Reference item created successfully")
}

func (ctrl *CmsController) UpdateReference(c *fiber.Ctx) error {
	if err := ctrl.verifyAdmin(c); err != nil {
		return response.Error(c, err)
	}

	refType := c.Params("type")
	code := c.Params("code")

	var req model.ReferenceRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request body", "INVALID_REQUEST_BODY")
	}

	result, err := ctrl.cmsService.UpdateReference(refType, code, &req)
	if err != nil {
		return response.Error(c, err)
	}
	return response.Success(c, result, "Reference item updated successfully")
}

func (ctrl *CmsController) DeleteReference(c *fiber.Ctx) error {
	if err := ctrl.verifyAdmin(c); err != nil {
		return response.Error(c, err)
	}

	refType := c.Params("type")
	code := c.Params("code")

	if err := ctrl.cmsService.DeleteReference(refType, code); err != nil {
		return response.Error(c, err)
	}
	return response.Success(c, nil, "Reference item deleted successfully")
}
