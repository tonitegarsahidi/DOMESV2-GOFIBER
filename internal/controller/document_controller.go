package controller

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"domesv2/internal/model"
	"domesv2/internal/service"
	"domesv2/pkg/response"
)

type DocumentController struct {
	docService service.DocumentService
}

func NewDocumentController(docService service.DocumentService) *DocumentController {
	return &DocumentController{
		docService: docService,
	}
}

func (ctrl *DocumentController) ListPublic(c *fiber.Ctx) error {
	filters := make(map[string]interface{})

	// Parse query params
	if pageStr := c.Query("page"); pageStr != "" {
		if val, err := strconv.Atoi(pageStr); err == nil {
			filters["page"] = val
		}
	}
	if limitStr := c.Query("limit"); limitStr != "" {
		if val, err := strconv.Atoi(limitStr); err == nil {
			filters["limit"] = val
		}
	}
	filters["sort"] = c.Query("sort")
	filters["agencies"] = c.Query("agencies")
	filters["sdgs"] = c.Query("sdgs")
	filters["sectors"] = c.Query("sectors")
	filters["langs"] = c.Query("langs")
	filters["jointProgrammes"] = c.Query("jointProgrammes")
	filters["lnobs"] = c.Query("lnobs")
	filters["nonUnPartners"] = c.Query("nonUnPartners")

	if yfStr := c.Query("yearFrom"); yfStr != "" {
		if val, err := strconv.Atoi(yfStr); err == nil {
			filters["yearFrom"] = val
		}
	}
	if ytStr := c.Query("yearTo"); ytStr != "" {
		if val, err := strconv.Atoi(ytStr); err == nil {
			filters["yearTo"] = val
		}
	}

	result, err := ctrl.docService.ListPublicDocuments(filters)
	if err != nil {
		return response.Error(c, err)
	}

	return response.Success(c, result, "Documents retrieved successfully")
}

func (ctrl *DocumentController) SearchPublic(c *fiber.Ctx) error {
	q := c.Query("q")
	sort := c.Query("sort", "relevance")
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 12)

	filters := make(map[string]interface{})
	filters["agencies"] = c.Query("agencies")
	filters["sdgs"] = c.Query("sdgs")
	filters["sectors"] = c.Query("sectors")
	filters["langs"] = c.Query("langs")
	filters["jointProgrammes"] = c.Query("jointProgrammes")
	filters["lnobs"] = c.Query("lnobs")
	filters["nonUnPartners"] = c.Query("nonUnPartners")

	if yfStr := c.Query("yearFrom"); yfStr != "" {
		if val, err := strconv.Atoi(yfStr); err == nil {
			filters["yearFrom"] = val
		}
	}
	if ytStr := c.Query("yearTo"); ytStr != "" {
		if val, err := strconv.Atoi(ytStr); err == nil {
			filters["yearTo"] = val
		}
	}

	result, err := ctrl.docService.SearchPublicDocuments(q, page, limit, sort, filters)
	if err != nil {
		return response.Error(c, err)
	}

	// The service already wraps the response matching search format
	return c.Status(fiber.StatusOK).JSON(result)
}

func (ctrl *DocumentController) GetByIDOrSlug(c *fiber.Ctx) error {
	idParam := c.Params("id")
	if idParam == "" {
		return response.BadRequest(c, "ID or Slug parameter is required", "VALIDATION_FAILED")
	}

	var result *model.DocumentResponse
	var err error

	// Detect if idParam is int
	if id, parseErr := strconv.Atoi(idParam); parseErr == nil {
		result, err = ctrl.docService.GetDocumentByID(uint(id))
	} else {
		result, err = ctrl.docService.GetDocumentBySlug(idParam)
	}

	if err != nil {
		return response.Error(c, err)
	}

	return response.Success(c, result, "Document retrieved successfully")
}

func (ctrl *DocumentController) GetRelated(c *fiber.Ctx) error {
	idParam := c.Params("id")
	var related []model.DocumentResponse
	var err error

	if id, parseErr := strconv.Atoi(idParam); parseErr == nil {
		related, err = ctrl.docService.GetRelatedDocuments(uint(id))
	} else {
		// Fetch doc by slug first to get ID
		docResp, getErr := ctrl.docService.GetDocumentBySlug(idParam)
		if getErr != nil {
			return response.Error(c, getErr)
		}
		related, err = ctrl.docService.GetRelatedDocuments(docResp.ID)
	}

	if err != nil {
		return response.Error(c, err)
	}

	return response.Success(c, related, "Related documents retrieved successfully")
}

func (ctrl *DocumentController) Download(c *fiber.Ctx) error {
	idParam := c.Params("id")
	var result map[string]interface{}
	var err error

	if id, parseErr := strconv.Atoi(idParam); parseErr == nil {
		result, err = ctrl.docService.GenerateDownloadLink(uint(id))
	} else {
		docResp, getErr := ctrl.docService.GetDocumentBySlug(idParam)
		if getErr != nil {
			return response.Error(c, getErr)
		}
		result, err = ctrl.docService.GenerateDownloadLink(docResp.ID)
	}

	if err != nil {
		return response.Error(c, err)
	}

	return response.Success(c, result, "Download link generated")
}

func (ctrl *DocumentController) GetPlatformStats(c *fiber.Ctx) error {
	result, err := ctrl.docService.GetPlatformStats()
	if err != nil {
		return response.Error(c, err)
	}
	return response.Success(c, result, "Platform statistics retrieved successfully")
}

func (ctrl *DocumentController) GetAnalyticsOverview(c *fiber.Ctx) error {
	result, err := ctrl.docService.GetAnalyticsOverview()
	if err != nil {
		return response.Error(c, err)
	}
	return response.Success(c, result, "Analytics overview retrieved successfully")
}

func (ctrl *DocumentController) GetUploadsOverTime(c *fiber.Ctx) error {
	fromYear := c.QueryInt("fromYear", 2014)
	toYear := c.QueryInt("toYear", 2024)

	result, err := ctrl.docService.GetUploadsOverTime(fromYear, toYear)
	if err != nil {
		return response.Error(c, err)
	}
	return response.Success(c, result, "Uploads over time retrieved successfully")
}

func (ctrl *DocumentController) GetBySdgAnalytics(c *fiber.Ctx) error {
	result, err := ctrl.docService.GetBySdgAnalytics()
	if err != nil {
		return response.Error(c, err)
	}
	return response.Success(c, result, "Documents by SDG retrieved successfully")
}

func (ctrl *DocumentController) GetByAgencyAnalytics(c *fiber.Ctx) error {
	result, err := ctrl.docService.GetByAgencyAnalytics()
	if err != nil {
		return response.Error(c, err)
	}
	return response.Success(c, result, "Documents by agency retrieved successfully")
}

func (ctrl *DocumentController) GetBySectorAnalytics(c *fiber.Ctx) error {
	result, err := ctrl.docService.GetBySectorAnalytics()
	if err != nil {
		return response.Error(c, err)
	}
	return response.Success(c, result, "Documents by sector retrieved successfully")
}

func (ctrl *DocumentController) GetByLanguageAnalytics(c *fiber.Ctx) error {
	result, err := ctrl.docService.GetByLanguageAnalytics()
	if err != nil {
		return response.Error(c, err)
	}
	return response.Success(c, result, "Documents by language retrieved successfully")
}

// CMS Endpoints
func (ctrl *DocumentController) ListSubmissions(c *fiber.Ctx) error {
	status := c.Query("status", "all")
	search := c.Query("search")
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 20)

	result, err := ctrl.docService.ListSubmissions(status, search, page, limit)
	if err != nil {
		return response.Error(c, err)
	}

	return response.Success(c, result, "Submissions retrieved successfully")
}

func (ctrl *DocumentController) CreateSubmission(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)

	var req model.SubmissionRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request body", "INVALID_REQUEST_BODY")
	}

	doc, err := ctrl.docService.CreateSubmission(userID, &req)
	if err != nil {
		return response.Error(c, err)
	}

	return response.Created(c, fiber.Map{
		"id":         doc.ID,
		"status":     doc.Status,
		"created_at": doc.CreatedAt,
	}, "Submission created successfully")
}

func (ctrl *DocumentController) SaveDraft(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)

	idParam := c.Params("id")
	var submissionID uint
	if id, err := strconv.Atoi(idParam); err == nil && idParam != "0" {
		submissionID = uint(id)
	}

	var req model.DraftRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request body", "INVALID_REQUEST_BODY")
	}

	doc, err := ctrl.docService.SaveDraft(userID, submissionID, req.Step, req.Data)
	if err != nil {
		return response.Error(c, err)
	}

	return response.Success(c, fiber.Map{
		"id":       doc.ID,
		"step":     req.Step,
		"saved_at": doc.UpdatedAt,
	}, "Draft saved successfully")
}

func (ctrl *DocumentController) DeleteSubmission(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return response.BadRequest(c, "Invalid submission ID", "VALIDATION_FAILED")
	}

	if err := ctrl.docService.DeleteSubmission(uint(id)); err != nil {
		return response.Error(c, err)
	}

	return response.Success(c, nil, "Submission deleted successfully")
}

func (ctrl *DocumentController) PublishDocument(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return response.BadRequest(c, "Invalid document ID", "VALIDATION_FAILED")
	}

	doc, err := ctrl.docService.PublishDocument(uint(id))
	if err != nil {
		return response.Error(c, err)
	}

	return response.Success(c, fiber.Map{
		"id":           doc.ID,
		"status":       doc.Status,
		"published_at": doc.UpdatedAt,
	}, "Document published successfully")
}

func (ctrl *DocumentController) UnpublishDocument(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return response.BadRequest(c, "Invalid document ID", "VALIDATION_FAILED")
	}

	doc, err := ctrl.docService.UnpublishDocument(uint(id))
	if err != nil {
		return response.Error(c, err)
	}

	return response.Success(c, fiber.Map{
		"id":     doc.ID,
		"status": doc.Status,
	}, "Document unpublished successfully")
}
