package repository

import (
	"strings"

	"domesv2/config/database"
	"domesv2/internal/model"
	"domesv2/pkg/errors"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type ReportRepository interface {
	Create(report *model.Report) error
	List(status string, search string, page int, limit int) ([]model.Report, int, error)
	UpdateStatus(id string, status string) (*model.Report, error)
}

type reportRepository struct {
	db *gorm.DB
}

func NewReportRepository() ReportRepository {
	return &reportRepository{
		db: database.GetDB(),
	}
}

func (r *reportRepository) Create(report *model.Report) error {
	if err := r.db.Create(report).Error; err != nil {
		zap.L().Error("Failed to create report", zap.Error(err))
		return errors.NewInternalServerError("Failed to submit report", "DATABASE_ERROR")
	}
	return nil
}

func (r *reportRepository) List(status string, search string, page int, limit int) ([]model.Report, int, error) {
	var reports []model.Report
	query := r.db.Model(&model.Report{}).Preload("Document")

	if status != "" && status != "all" {
		query = query.Where("status = ?", status)
	}

	if search != "" {
		// Join documents to search by document title
		query = query.Joins("JOIN V2Documents ON V2Documents.id = V2Reports.document_id").
			Where("LOWER(V2Documents.title) LIKE ?", "%"+strings.ToLower(search)+"%")
	}

	var totalItems int64
	query.Count(&totalItems)

	offset := (page - 1) * limit
	err := query.Order("createdAt desc").Limit(limit).Offset(offset).Find(&reports).Error
	if err != nil {
		zap.L().Error("Failed to fetch reports", zap.Error(err))
		return nil, 0, errors.NewInternalServerError("Failed to fetch reports", "DATABASE_ERROR")
	}

	return reports, int(totalItems), nil
}

func (r *reportRepository) UpdateStatus(id string, status string) (*model.Report, error) {
	var report model.Report
	if err := r.db.First(&report, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NewNotFoundError("Report not found", "REPORT_NOT_FOUND")
		}
		return nil, errors.NewInternalServerError("Database lookup error", "DATABASE_ERROR")
	}

	report.Status = status
	if err := r.db.Save(&report).Error; err != nil {
		zap.L().Error("Failed to update report status", zap.Error(err))
		return nil, errors.NewInternalServerError("Failed to update status", "DATABASE_ERROR")
	}

	return &report, nil
}
