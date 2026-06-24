package repository

import (
	"strings"

	"domesv2/config/database"
	"domesv2/internal/model"
	"domesv2/pkg/errors"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type CmsRepository interface {
	ListUsers(search string, role string, status string, page int, limit int) ([]model.User, int, error)
	GetTopDownloads(limit int) ([]model.Document, error)
	GetTopViews(limit int) ([]model.Document, error)
	GetDashboardStats() (map[string]interface{}, error)
}

type cmsRepository struct {
	db *gorm.DB
}

func NewCmsRepository() CmsRepository {
	return &cmsRepository{
		db: database.GetDB(),
	}
}

func (r *cmsRepository) ListUsers(search string, role string, status string, page int, limit int) ([]model.User, int, error) {
	var users []model.User
	query := r.db.Model(&model.User{})

	if role != "" {
		query = query.Where("role = ?", role)
	}

	if status != "" {
		query = query.Where("status = ?", status)
	}

	if search != "" {
		searchText := "%" + strings.ToLower(search) + "%"
		query = query.Where("LOWER(email) LIKE ? OR LOWER(first_name) LIKE ? OR LOWER(last_name) LIKE ? OR LOWER(organization) LIKE ?",
			searchText, searchText, searchText, searchText)
	}

	var totalItems int64
	query.Count(&totalItems)

	offset := (page - 1) * limit
	err := query.Order("createdAt desc").Limit(limit).Offset(offset).Find(&users).Error
	if err != nil {
		zap.L().Error("Failed to fetch users", zap.Error(err))
		return nil, 0, errors.NewInternalServerError("Failed to fetch users", "DATABASE_ERROR")
	}

	return users, int(totalItems), nil
}

func (r *cmsRepository) GetTopDownloads(limit int) ([]model.Document, error) {
	var docs []model.Document
	err := r.db.Model(&model.Document{}).Where("status = ?", "published").
		Order("downloads desc, title asc").Limit(limit).Find(&docs).Error
	if err != nil {
		zap.L().Error("Failed to fetch top downloads", zap.Error(err))
		return nil, errors.NewInternalServerError("Failed to fetch top downloads", "DATABASE_ERROR")
	}
	return docs, nil
}

func (r *cmsRepository) GetTopViews(limit int) ([]model.Document, error) {
	var docs []model.Document
	err := r.db.Model(&model.Document{}).Where("status = ?", "published").
		Order("views desc, title asc").Limit(limit).Find(&docs).Error
	if err != nil {
		zap.L().Error("Failed to fetch top views", zap.Error(err))
		return nil, errors.NewInternalServerError("Failed to fetch top views", "DATABASE_ERROR")
	}
	return docs, nil
}

func (r *cmsRepository) GetDashboardStats() (map[string]interface{}, error) {
	var totalDocs int64
	var totalUsers int64
	var pendingApprovals int64
	var pendingReports int64

	// GORM count queries
	r.db.Model(&model.Document{}).Count(&totalDocs)
	r.db.Model(&model.User{}).Count(&totalUsers)
	r.db.Model(&model.Document{}).Where("status = ?", "pending_review").Count(&pendingApprovals)
	r.db.Model(&model.Report{}).Where("status = ?", "open").Count(&pendingReports)

	var stats struct {
		Views     int64
		Downloads int64
	}
	r.db.Model(&model.Document{}).Select("SUM(views) as views, SUM(downloads) as downloads").Scan(&stats)

	return map[string]interface{}{
		"total_documents":   int(totalDocs),
		"total_users":       int(totalUsers),
		"pending_approvals": int(pendingApprovals),
		"reports":           int(pendingReports),
		"total_views":       int(stats.Views),
		"total_downloads":   int(stats.Downloads),
	}, nil
}
