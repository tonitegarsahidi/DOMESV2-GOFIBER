package service

import (
	"domesv2/internal/model"
	"domesv2/internal/repository"
	"domesv2/pkg/captcha"
	"domesv2/pkg/errors"
)

type ReportService interface {
	SubmitReport(req *model.CreateReportRequest) (*model.Report, error)
	ListReports(status string, search string, page int, limit int) (map[string]interface{}, error)
	UpdateReportStatus(id string, status string) (*model.Report, error)
}

type reportService struct {
	reportRepo repository.ReportRepository
}

func NewReportService(reportRepo repository.ReportRepository) ReportService {
	return &reportService{
		reportRepo: reportRepo,
	}
}

func (s *reportService) SubmitReport(req *model.CreateReportRequest) (*model.Report, error) {
	if err := captcha.VerifyCaptcha(req.Captcha); err != nil {
		return nil, err
	}

	if req.DocumentID == "" {
		return nil, errors.NewValidationError("Document ID is required", "VALIDATION_FAILED")
	}
	if req.ReporterName == "" || req.ReporterEmail == "" || req.Details == "" {
		return nil, errors.NewValidationError("Reporter name, email, and details are required", "VALIDATION_FAILED")
	}

	report := &model.Report{
		DocumentID:    req.DocumentID,
		ReporterName:  req.ReporterName,
		ReporterEmail: req.ReporterEmail,
		Details:       req.Details,
		Status:        "open",
	}

	if err := s.reportRepo.Create(report); err != nil {
		return nil, err
	}

	return report, nil
}

func (s *reportService) ListReports(status string, search string, page int, limit int) (map[string]interface{}, error) {
	reports, totalItems, err := s.reportRepo.List(status, search, page, limit)
	if err != nil {
		return nil, err
	}

	var items []map[string]interface{}
	for _, r := range reports {
		docTitle := ""
		if r.Document != nil {
			docTitle = r.Document.Title
		}

		items = append(items, map[string]interface{}{
			"id":             r.ID,
			"document_id":    r.DocumentID,
			"document_title": docTitle,
			"reporter_name":  r.ReporterName,
			"reporter_email": r.ReporterEmail,
			"details":        r.Details,
			"status":         r.Status,
			"created_at":     r.CreatedAt,
		})
	}

	totalPages := 0
	if totalItems > 0 {
		totalPages = (totalItems + limit - 1) / limit
	}

	return map[string]interface{}{
		"items": items,
		"pagination": map[string]interface{}{
			"page":       page,
			"limit":      limit,
			"totalItems": totalItems,
			"totalPages": totalPages,
		},
	}, nil
}

func (s *reportService) UpdateReportStatus(id string, status string) (*model.Report, error) {
	if status != "open" && status != "in_progress" && status != "resolved" {
		return nil, errors.NewValidationError("Status must be one of: open, in_progress, resolved", "VALIDATION_FAILED")
	}

	return s.reportRepo.UpdateStatus(id, status)
}
