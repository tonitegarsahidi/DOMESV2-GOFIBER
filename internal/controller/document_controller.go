package controller

import (
	"strconv"
	"strings"

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

	if q := c.Query("q"); q != "" {
		filters["q"] = q
	}
	if agencies := c.Query("agencies"); agencies != "" {
		filters["agencies"] = agencies
	}
	if jointProgrammes := c.Query("jointProgrammes"); jointProgrammes != "" {
		filters["jointProgrammes"] = jointProgrammes
	}
	if yearFromStr := c.Query("yearFrom"); yearFromStr != "" {
		if val, err := strconv.Atoi(yearFromStr); err == nil {
			filters["yearFrom"] = val
		}
	}
	if yearToStr := c.Query("yearTo"); yearToStr != "" {
		if val, err := strconv.Atoi(yearToStr); err == nil {
			filters["yearTo"] = val
		}
	}
	if langs := c.Query("langs"); langs != "" {
		filters["langs"] = langs
	}
	if sdgs := c.Query("sdgs"); sdgs != "" {
		filters["sdgs"] = sdgs
	}
	if sectors := c.Query("sectors"); sectors != "" {
		filters["sectors"] = sectors
	}
	if lnobs := c.Query("lnobs"); lnobs != "" {
		filters["lnobs"] = lnobs
	}
	if nonUnPartners := c.Query("nonUnPartners"); nonUnPartners != "" {
		filters["nonUnPartners"] = nonUnPartners
	}
	if sort := c.Query("sort"); sort != "" {
		filters["sort"] = sort
	}
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

	result, err := ctrl.docService.ListPublicDocuments(filters)
	if err != nil {
		return response.Error(c, err)
	}

	return response.Success(c, result, "Documents retrieved successfully")
}

func (ctrl *DocumentController) SearchPublic(c *fiber.Ctx) error {
	q := c.Query("q")
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 12)
	sort := c.Query("sort", "newest")

	filters := make(map[string]interface{})
	if agencies := c.Query("agencies"); agencies != "" {
		filters["agencies"] = agencies
	}
	if jointProgrammes := c.Query("jointProgrammes"); jointProgrammes != "" {
		filters["jointProgrammes"] = jointProgrammes
	}
	if yearFromStr := c.Query("yearFrom"); yearFromStr != "" {
		if val, err := strconv.Atoi(yearFromStr); err == nil {
			filters["yearFrom"] = val
		}
	}
	if yearToStr := c.Query("yearTo"); yearToStr != "" {
		if val, err := strconv.Atoi(yearToStr); err == nil {
			filters["yearTo"] = val
		}
	}
	if langs := c.Query("langs"); langs != "" {
		filters["langs"] = langs
	}
	if sdgs := c.Query("sdgs"); sdgs != "" {
		filters["sdgs"] = sdgs
	}
	if sectors := c.Query("sectors"); sectors != "" {
		filters["sectors"] = sectors
	}
	if lnobs := c.Query("lnobs"); lnobs != "" {
		filters["lnobs"] = lnobs
	}
	if nonUnPartners := c.Query("nonUnPartners"); nonUnPartners != "" {
		filters["nonUnPartners"] = nonUnPartners
	}

	result, err := ctrl.docService.SearchPublicDocuments(q, page, limit, sort, filters)
	if err != nil {
		return response.Error(c, err)
	}

	return c.JSON(result)
}

func (ctrl *DocumentController) GetByIDOrSlug(c *fiber.Ctx) error {
	idParam := c.Params("id")
	if idParam == "" {
		return response.BadRequest(c, "ID or Slug parameter is required", "VALIDATION_FAILED")
	}

	var result *model.DocumentResponse
	var err error

	if len(idParam) == 36 && strings.Contains(idParam, "-") {
		result, err = ctrl.docService.GetDocumentByID(idParam)
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

	if len(idParam) == 36 && strings.Contains(idParam, "-") {
		related, err = ctrl.docService.GetRelatedDocuments(idParam)
	} else {
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

	if len(idParam) == 36 && strings.Contains(idParam, "-") {
		result, err = ctrl.docService.GenerateDownloadLink(idParam)
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

	return response.Success(c, result, "Download link generated successfully")
}

func (ctrl *DocumentController) GetPlatformStats(c *fiber.Ctx) error {
	result, err := ctrl.docService.GetPlatformStats()
	if err != nil {
		return response.Error(c, err)
	}
	return response.Success(c, result, "Platform stats retrieved successfully")
}

func (ctrl *DocumentController) GetAnalyticsOverview(c *fiber.Ctx) error {
	result, err := ctrl.docService.GetAnalyticsOverview()
	if err != nil {
		return response.Error(c, err)
	}
	return response.Success(c, result, "Analytics overview retrieved successfully")
}

func (ctrl *DocumentController) GetUploadsOverTime(c *fiber.Ctx) error {
	fromYear := c.QueryInt("fromYear", 2019)
	toYear := c.QueryInt("toYear", 2024)

	result, err := ctrl.docService.GetUploadsOverTime(fromYear, toYear)
	if err != nil {
		return response.Error(c, err)
	}
	return response.Success(c, result, "Uploads over time analytics retrieved successfully")
}

func (ctrl *DocumentController) GetBySdgAnalytics(c *fiber.Ctx) error {
	result, err := ctrl.docService.GetBySdgAnalytics()
	if err != nil {
		return response.Error(c, err)
	}
	return response.Success(c, result, "Analytics by SDG retrieved successfully")
}

func (ctrl *DocumentController) GetByAgencyAnalytics(c *fiber.Ctx) error {
	result, err := ctrl.docService.GetByAgencyAnalytics()
	if err != nil {
		return response.Error(c, err)
	}
	return response.Success(c, result, "Analytics by agency retrieved successfully")
}

func (ctrl *DocumentController) GetBySectorAnalytics(c *fiber.Ctx) error {
	result, err := ctrl.docService.GetBySectorAnalytics()
	if err != nil {
		return response.Error(c, err)
	}
	return response.Success(c, result, "Analytics by sector retrieved successfully")
}

func (ctrl *DocumentController) GetByLanguageAnalytics(c *fiber.Ctx) error {
	result, err := ctrl.docService.GetByLanguageAnalytics()
	if err != nil {
		return response.Error(c, err)
	}
	return response.Success(c, result, "Analytics by language retrieved successfully")
}

func (ctrl *DocumentController) ListSubmissions(c *fiber.Ctx) error {
	status := c.Query("status", "all")
	search := c.Query("search")
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)

	result, err := ctrl.docService.ListSubmissions(status, search, page, limit)
	if err != nil {
		return response.Error(c, err)
	}

	return response.Success(c, result, "Submissions list retrieved successfully")
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
		"code":       doc.Code,
		"slug":       doc.Slug,
		"title":      doc.Title,
		"status":     doc.Status,
		"created_at": doc.CreatedAt,
	}, "Submission created successfully")
}

func (ctrl *DocumentController) SaveDraft(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)

	idParam := c.Params("id")
	var submissionID string
	if idParam != "" && idParam != "0" {
		submissionID = idParam
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
	if idParam == "" {
		return response.BadRequest(c, "Invalid submission ID", "VALIDATION_FAILED")
	}

	if err := ctrl.docService.DeleteSubmission(idParam); err != nil {
		return response.Error(c, err)
	}

	return response.Success(c, nil, "Submission deleted successfully")
}

func (ctrl *DocumentController) PublishDocument(c *fiber.Ctx) error {
	idParam := c.Params("id")
	if idParam == "" {
		return response.BadRequest(c, "Invalid document ID", "VALIDATION_FAILED")
	}

	doc, err := ctrl.docService.PublishDocument(idParam)
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
	if idParam == "" {
		return response.BadRequest(c, "Invalid document ID", "VALIDATION_FAILED")
	}

	doc, err := ctrl.docService.UnpublishDocument(idParam)
	if err != nil {
		return response.Error(c, err)
	}

	return response.Success(c, fiber.Map{
		"id":     doc.ID,
		"status": doc.Status,
	}, "Document unpublished successfully")
}
