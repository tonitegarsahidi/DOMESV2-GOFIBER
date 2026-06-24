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
	ListReferences(refType string) (interface{}, error)
	GetReferenceByCode(refType string, code string) (interface{}, error)
	CreateReference(refType string, item interface{}) error
	UpdateReference(refType string, code string, item interface{}) error
	DeleteReference(refType string, code string) error
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

func (r *cmsRepository) ListReferences(refType string) (interface{}, error) {
	switch strings.ToLower(refType) {
	case "agencies":
		var list []model.Agency
		if err := r.db.Order("code asc").Find(&list).Error; err != nil {
			return nil, errors.NewInternalServerError("Failed to fetch agencies", "DATABASE_ERROR")
		}
		return list, nil
	case "sdgs":
		var list []model.Sdg
		if err := r.db.Order("CAST(SUBSTRING(code, 6) AS UNSIGNED) asc").Find(&list).Error; err != nil {
			return nil, errors.NewInternalServerError("Failed to fetch SDGs", "DATABASE_ERROR")
		}
		return list, nil
	case "sectors":
		var list []model.Sector
		if err := r.db.Order("name asc").Find(&list).Error; err != nil {
			return nil, errors.NewInternalServerError("Failed to fetch sectors", "DATABASE_ERROR")
		}
		return list, nil
	case "languages":
		var list []model.Language
		if err := r.db.Order("name asc").Find(&list).Error; err != nil {
			return nil, errors.NewInternalServerError("Failed to fetch languages", "DATABASE_ERROR")
		}
		return list, nil
	case "joint-programmes":
		var list []model.JointProgramme
		if err := r.db.Order("name asc").Find(&list).Error; err != nil {
			return nil, errors.NewInternalServerError("Failed to fetch joint programmes", "DATABASE_ERROR")
		}
		return list, nil
	case "lnobs":
		var list []model.Lnob
		if err := r.db.Order("name asc").Find(&list).Error; err != nil {
			return nil, errors.NewInternalServerError("Failed to fetch LNOBs", "DATABASE_ERROR")
		}
		return list, nil
	case "non-un-partners":
		var list []model.NonUnPartner
		if err := r.db.Order("name asc").Find(&list).Error; err != nil {
			return nil, errors.NewInternalServerError("Failed to fetch non-UN partners", "DATABASE_ERROR")
		}
		return list, nil
	case "organizations":
		var list []model.Organization
		if err := r.db.Order("name asc").Find(&list).Error; err != nil {
			return nil, errors.NewInternalServerError("Failed to fetch organizations", "DATABASE_ERROR")
		}
		return list, nil
	default:
		return nil, errors.NewValidationError("Invalid reference type", "INVALID_REF_TYPE")
	}
}

func (r *cmsRepository) GetReferenceByCode(refType string, code string) (interface{}, error) {
	switch strings.ToLower(refType) {
	case "agencies":
		var item model.Agency
		if err := r.db.Where("code = ?", code).First(&item).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return nil, errors.NewNotFoundError("Agency not found", "REF_NOT_FOUND")
			}
			return nil, errors.NewInternalServerError("Database error", "DATABASE_ERROR")
		}
		return &item, nil
	case "sdgs":
		var item model.Sdg
		if err := r.db.Where("code = ?", code).First(&item).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return nil, errors.NewNotFoundError("SDG not found", "REF_NOT_FOUND")
			}
			return nil, errors.NewInternalServerError("Database error", "DATABASE_ERROR")
		}
		return &item, nil
	case "sectors":
		var item model.Sector
		if err := r.db.Where("code = ?", code).First(&item).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return nil, errors.NewNotFoundError("Sector not found", "REF_NOT_FOUND")
			}
			return nil, errors.NewInternalServerError("Database error", "DATABASE_ERROR")
		}
		return &item, nil
	case "languages":
		var item model.Language
		if err := r.db.Where("code = ?", code).First(&item).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return nil, errors.NewNotFoundError("Language not found", "REF_NOT_FOUND")
			}
			return nil, errors.NewInternalServerError("Database error", "DATABASE_ERROR")
		}
		return &item, nil
	case "joint-programmes":
		var item model.JointProgramme
		if err := r.db.Where("code = ?", code).First(&item).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return nil, errors.NewNotFoundError("Joint programme not found", "REF_NOT_FOUND")
			}
			return nil, errors.NewInternalServerError("Database error", "DATABASE_ERROR")
		}
		return &item, nil
	case "lnobs":
		var item model.Lnob
		if err := r.db.Where("code = ?", code).First(&item).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return nil, errors.NewNotFoundError("LNOB group not found", "REF_NOT_FOUND")
			}
			return nil, errors.NewInternalServerError("Database error", "DATABASE_ERROR")
		}
		return &item, nil
	case "non-un-partners":
		var item model.NonUnPartner
		if err := r.db.Where("code = ?", code).First(&item).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return nil, errors.NewNotFoundError("Non-UN partner not found", "REF_NOT_FOUND")
			}
			return nil, errors.NewInternalServerError("Database error", "DATABASE_ERROR")
		}
		return &item, nil
	case "organizations":
		var item model.Organization
		if err := r.db.Where("code = ?", code).First(&item).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return nil, errors.NewNotFoundError("Organization not found", "REF_NOT_FOUND")
			}
			return nil, errors.NewInternalServerError("Database error", "DATABASE_ERROR")
		}
		return &item, nil
	default:
		return nil, errors.NewValidationError("Invalid reference type", "INVALID_REF_TYPE")
	}
}

func (r *cmsRepository) CreateReference(refType string, item interface{}) error {
	var code string
	switch v := item.(type) {
	case *model.Agency:
		code = v.Code
	case *model.Sdg:
		code = v.Code
	case *model.Sector:
		code = v.Code
	case *model.Language:
		code = v.Code
	case *model.JointProgramme:
		code = v.Code
	case *model.Lnob:
		code = v.Code
	case *model.NonUnPartner:
		code = v.Code
	case *model.Organization:
		code = v.Code
	}

	if code == "" {
		return errors.NewValidationError("Code is required", "VALIDATION_FAILED")
	}

	existing, _ := r.GetReferenceByCode(refType, code)
	if existing != nil {
		return errors.NewConflictError("Reference item with this code already exists", "REF_CODE_EXISTS")
	}

	if err := r.db.Create(item).Error; err != nil {
		zap.L().Error("Failed to create reference item", zap.Error(err))
		return errors.NewInternalServerError("Failed to save reference item", "DATABASE_ERROR")
	}
	return nil
}

func (r *cmsRepository) UpdateReference(refType string, code string, item interface{}) error {
	_, err := r.GetReferenceByCode(refType, code)
	if err != nil {
		return err
	}

	if err := r.db.Save(item).Error; err != nil {
		zap.L().Error("Failed to update reference item", zap.Error(err))
		return errors.NewInternalServerError("Failed to update reference item", "DATABASE_ERROR")
	}
	return nil
}

func (r *cmsRepository) DeleteReference(refType string, code string) error {
	existing, err := r.GetReferenceByCode(refType, code)
	if err != nil {
		return err
	}

	if err := r.db.Delete(existing).Error; err != nil {
		zap.L().Error("Failed to delete reference item", zap.Error(err))
		return errors.NewInternalServerError("Failed to delete reference item", "DATABASE_ERROR")
	}
	return nil
}
