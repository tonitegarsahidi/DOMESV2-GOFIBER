package controller

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"domesv2/internal/model"
	"domesv2/internal/service"
	"domesv2/pkg/response"
)

type ReportController struct {
	reportService service.ReportService
}

func NewReportController(reportService service.ReportService) *ReportController {
	return &ReportController{
		reportService: reportService,
	}
}

func (ctrl *ReportController) SubmitReport(c *fiber.Ctx) error {
	var req model.CreateReportRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request body", "INVALID_REQUEST_BODY")
	}

	report, err := ctrl.reportService.SubmitReport(&req)
	if err != nil {
		return response.Error(c, err)
	}

	return response.Created(c, fiber.Map{
		"id":          report.ID,
		"document_id": report.DocumentID,
		"status":      report.Status,
		"created_at":  report.CreatedAt,
	}, "Report submitted successfully")
}

func (ctrl *ReportController) ListReports(c *fiber.Ctx) error {
	status := c.Query("status", "all")
	search := c.Query("search")
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 20)

	result, err := ctrl.reportService.ListReports(status, search, page, limit)
	if err != nil {
		return response.Error(c, err)
	}

	return response.Success(c, result, "Reports retrieved successfully")
}

func (ctrl *ReportController) UpdateStatus(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return response.BadRequest(c, "Invalid report ID", "VALIDATION_FAILED")
	}

	var req model.UpdateReportStatusRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request body", "INVALID_REQUEST_BODY")
	}

	report, err := ctrl.reportService.UpdateReportStatus(uint(id), req.Status)
	if err != nil {
		return response.Error(c, err)
	}

	return response.Success(c, fiber.Map{
		"id":         report.ID,
		"status":     report.Status,
		"updated_at": report.UpdatedAt,
	}, "Report status updated successfully")
}
