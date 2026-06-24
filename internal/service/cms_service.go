package service

import (
	"strings"
	"time"

	"domesv2/config/database"
	"domesv2/internal/model"
	"domesv2/internal/repository"
	"domesv2/pkg/errors"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type CmsService interface {
	GetDashboardStats() (map[string]interface{}, error)
	GetRecentActivity(limit int) ([]map[string]interface{}, error)
	GetAnalyticsSummary(period string) (map[string]interface{}, error)
	GetTopDownloads(limit int) ([]map[string]interface{}, error)
	GetTopViews(limit int) ([]map[string]interface{}, error)
	ListUsers(search string, role string, status string, page int, limit int) (map[string]interface{}, error)
	CreateUser(req *model.CreateUserRequest) (*model.User, error)
	UpdateUser(id uint, req *model.UpdateUserRequest) (*model.User, error)
	DeleteUser(id uint) error
	ListReferences(refType string) (interface{}, error)
	CreateReference(refType string, req *model.ReferenceRequest) (interface{}, error)
	UpdateReference(refType string, code string, req *model.ReferenceRequest) (interface{}, error)
	DeleteReference(refType string, code string) error
}

type cmsService struct {
	cmsRepo  repository.CmsRepository
	userRepo repository.UserRepository
}

func NewCmsService(cmsRepo repository.CmsRepository, userRepo repository.UserRepository) CmsService {
	return &cmsService{
		cmsRepo:  cmsRepo,
		userRepo: userRepo,
	}
}

func (s *cmsService) GetDashboardStats() (map[string]interface{}, error) {
	rawStats, err := s.cmsRepo.GetDashboardStats()
	if err != nil {
		return nil, err
	}

	// Format data structure matching contract
	return map[string]interface{}{
		"stats": map[string]interface{}{
			"total_documents": map[string]interface{}{
				"value": rawStats["total_documents"],
				"change": 12.5,
				"trend": "up",
			},
			"total_views": map[string]interface{}{
				"value": rawStats["total_views"],
				"change": 8.3,
				"trend": "up",
			},
			"total_downloads": map[string]interface{}{
				"value": rawStats["total_downloads"],
				"change": -2.1,
				"trend": "down",
			},
			"total_users": map[string]interface{}{
				"value": rawStats["total_users"],
				"change": 5.7,
				"trend": "up",
			},
			"pending_approvals": map[string]interface{}{
				"value": rawStats["pending_approvals"],
				"change": 0,
				"trend": "neutral",
			},
			"reports": map[string]interface{}{
				"value": rawStats["reports"],
				"change": -1,
				"trend": "down",
			},
		},
	}, nil
}

func (s *cmsService) GetRecentActivity(limit int) ([]map[string]interface{}, error) {
	// Let's create an activity logger or mock list that is highly readable
	now := time.Now()
	activity := []map[string]interface{}{
		{
			"id":          1,
			"type":        "submission",
			"action":      "created",
			"description": "New document submitted: 'Digital Economy Report 2024'",
			"user":        "Erlangga Agustino",
			"user_avatar": "/uploads/avatars/erlangga.jpg",
			"timestamp":   now.Add(-2 * time.Minute).Format(time.RFC3339),
			"time_ago":    "2 minutes ago",
		},
		{
			"id":          2,
			"type":        "approval",
			"action":      "approved",
			"description": "Document 'Climate Change Adaptation' has been approved",
			"user":        "Admin User",
			"user_avatar": nil,
			"timestamp":   now.Add(-17 * time.Minute).Format(time.RFC3339),
			"time_ago":    "17 minutes ago",
		},
		{
			"id":          3,
			"type":        "report",
			"action":      "created",
			"description": "Broken link reported on document 'Water Sanitation Report'",
			"user":        "Budi Santoso",
			"user_avatar": nil,
			"timestamp":   now.Add(-32 * time.Minute).Format(time.RFC3339),
			"time_ago":    "32 minutes ago",
		},
	}

	if limit < len(activity) {
		return activity[:limit], nil
	}
	return activity, nil
}

func (s *cmsService) GetAnalyticsSummary(period string) (map[string]interface{}, error) {
	// Retrieve period totals (Mocking trend changes as per period)
	rawStats, err := s.cmsRepo.GetDashboardStats()
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"total_downloads": map[string]interface{}{
			"value":  rawStats["total_downloads"],
			"change": 12.5,
			"trend":  "up",
		},
		"total_views": map[string]interface{}{
			"value":  rawStats["total_views"],
			"change": 8.2,
			"trend":  "up",
		},
		"active_users": map[string]interface{}{
			"value":  3240, // standard mockup active user metric
			"change": -2.1,
			"trend":  "down",
		},
	}, nil
}

func (s *cmsService) GetTopDownloads(limit int) ([]map[string]interface{}, error) {
	docs, err := s.cmsRepo.GetTopDownloads(limit)
	if err != nil {
		return nil, err
	}

	var list []map[string]interface{}
	maxDownloads := 1
	if len(docs) > 0 && docs[0].Downloads > 0 {
		maxDownloads = docs[0].Downloads
	}

	for _, d := range docs {
		progress := (d.Downloads * 100) / maxDownloads
		list = append(list, map[string]interface{}{
			"title":     d.Title,
			"downloads": d.Downloads,
			"progress":  progress,
		})
	}
	return list, nil
}

func (s *cmsService) GetTopViews(limit int) ([]map[string]interface{}, error) {
	docs, err := s.cmsRepo.GetTopViews(limit)
	if err != nil {
		return nil, err
	}

	var list []map[string]interface{}
	for _, d := range docs {
		list = append(list, map[string]interface{}{
			"title":    d.Title,
			"category": d.LeadAgencyCode,
			"views":    d.Views,
		})
	}
	return list, nil
}

func (s *cmsService) ListUsers(search string, role string, status string, page int, limit int) (map[string]interface{}, error) {
	users, totalItems, err := s.cmsRepo.ListUsers(search, role, status, page, limit)
	if err != nil {
		return nil, err
	}

	var items []map[string]interface{}
	for _, u := range users {
		emailVerified := u.CreatedAt
		items = append(items, map[string]interface{}{
			"id":           u.ID,
			"first_name":   u.FirstName,
			"last_name":    u.LastName,
			"email":        u.Email,
			"phone_number": u.PhoneNumber,
			"organization": u.Organization,
			"position":     u.Position,
			"role":         u.Role,
			"status":       u.Status,
			"avatar_url":   u.AvatarURL,
			"created_at":   u.CreatedAt,
			"last_login":   emailVerified, // fallback to createdAt for mockup API
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

func (s *cmsService) CreateUser(req *model.CreateUserRequest) (*model.User, error) {
	if req.Email == "" || req.FirstName == "" || req.LastName == "" {
		return nil, errors.NewValidationError("First name, last name, and email are required", "VALIDATION_FAILED")
	}
	if req.Password == "" || len(req.Password) < 6 {
		return nil, errors.NewValidationError("Password must be at least 6 characters", "VALIDATION_FAILED")
	}
	if req.Password != req.ConfirmPassword {
		return nil, errors.NewValidationError("Passwords do not match", "VALIDATION_FAILED")
	}

	// Check if user already exists
	existing, _ := s.userRepo.FindByEmail(req.Email)
	if existing != nil {
		return nil, errors.NewConflictError("User with this email already exists", "USER_ALREADY_EXISTS")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		zap.L().Error("Failed to hash user password in CMS", zap.Error(err))
		return nil, errors.NewInternalServerError("Failed to create user", "PASSWORD_HASH_ERROR")
	}

	fullName := req.FirstName + " " + req.LastName
	username := strings.Split(req.Email, "@")[0]

	userRole := req.Role
	if userRole == "" {
		userRole = "editor"
	}
	userStatus := req.Status
	if userStatus == "" {
		userStatus = "active"
	}

	userType := "user"
	if strings.ToLower(userRole) == "administrator" {
		userType = "admin"
	}

	user := &model.User{
		Username:     &username,
		Name:         &fullName,
		FirstName:    &req.FirstName,
		LastName:     &req.LastName,
		Email:        req.Email,
		Password:     string(hashedPassword),
		Organization: &req.Organization,
		Position:     &req.Position,
		PhoneNumber:  &req.PhoneNumber,
		Role:         &userRole,
		Status:       &userStatus,
		Type:         &userType,
		NotificationPreferences: &model.NotificationPreference{
			DocumentApprovals:  true,
			BrokenLinkReports:  true,
			SystemUpdates:      false,
			EmailNotifications: true,
		},
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *cmsService) UpdateUser(id uint, req *model.UpdateUserRequest) (*model.User, error) {
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	if req.FirstName != nil {
		user.FirstName = req.FirstName
	}
	if req.LastName != nil {
		user.LastName = req.LastName
	}
	if req.Organization != nil {
		user.Organization = req.Organization
	}
	if req.Position != nil {
		user.Position = req.Position
	}
	if req.PhoneNumber != nil {
		user.PhoneNumber = req.PhoneNumber
	}
	if req.Role != nil {
		user.Role = req.Role
		userType := "user"
		if strings.ToLower(*req.Role) == "administrator" {
			userType = "admin"
		}
		user.Type = &userType
	}
	if req.Status != nil {
		user.Status = req.Status
	}

	// Update name if first_name or last_name changed
	fullName := *user.FirstName + " " + *user.LastName
	user.Name = &fullName

	if err := s.userRepo.Update(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *cmsService) DeleteUser(id uint) error {
	// Check if user exists
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return err
	}

	// Clear notification preferences first
	if user.NotificationPreferences != nil {
		// Handled via GORM delete constraint cascade, or explicitly
	}

	// Delete user directly
	db := database.GetDB()
	if err := db.Delete(user).Error; err != nil {
		zap.L().Error("Failed to delete user", zap.Error(err))
		return errors.NewInternalServerError("Failed to delete user", "DATABASE_ERROR")
	}

	return nil
}

func (s *cmsService) ListReferences(refType string) (interface{}, error) {
	return s.cmsRepo.ListReferences(refType)
}

func (s *cmsService) CreateReference(refType string, req *model.ReferenceRequest) (interface{}, error) {
	if req.Code == "" || req.Name == "" {
		return nil, errors.NewValidationError("Code and name are required", "VALIDATION_FAILED")
	}

	var item interface{}
	switch strings.ToLower(refType) {
	case "agencies":
		item = &model.Agency{Code: req.Code, Name: req.Name, LogoURL: req.LogoURL}
	case "sdgs":
		item = &model.Sdg{Code: req.Code, Name: req.Name, Icon: req.Icon, Color: req.Color}
	case "sectors":
		item = &model.Sector{Code: req.Code, Name: req.Name}
	case "languages":
		item = &model.Language{Code: req.Code, Name: req.Name}
	case "joint-programmes":
		item = &model.JointProgramme{Code: req.Code, Name: req.Name}
	case "lnobs":
		item = &model.Lnob{Code: req.Code, Name: req.Name}
	case "non-un-partners":
		item = &model.NonUnPartner{Code: req.Code, Name: req.Name}
	case "organizations":
		item = &model.Organization{Code: req.Code, Name: req.Name}
	default:
		return nil, errors.NewValidationError("Invalid reference type", "INVALID_REF_TYPE")
	}

	if err := s.cmsRepo.CreateReference(refType, item); err != nil {
		return nil, err
	}
	return item, nil
}

func (s *cmsService) UpdateReference(refType string, code string, req *model.ReferenceRequest) (interface{}, error) {
	if req.Name == "" {
		return nil, errors.NewValidationError("Name is required", "VALIDATION_FAILED")
	}

	existing, err := s.cmsRepo.GetReferenceByCode(refType, code)
	if err != nil {
		return nil, err
	}

	switch v := existing.(type) {
	case *model.Agency:
		v.Name = req.Name
		if req.LogoURL != "" {
			v.LogoURL = req.LogoURL
		}
	case *model.Sdg:
		v.Name = req.Name
		if req.Icon != "" {
			v.Icon = req.Icon
		}
		if req.Color != "" {
			v.Color = req.Color
		}
	case *model.Sector:
		v.Name = req.Name
	case *model.Language:
		v.Name = req.Name
	case *model.JointProgramme:
		v.Name = req.Name
	case *model.Lnob:
		v.Name = req.Name
	case *model.NonUnPartner:
		v.Name = req.Name
	case *model.Organization:
		v.Name = req.Name
	}

	if err := s.cmsRepo.UpdateReference(refType, code, existing); err != nil {
		return nil, err
	}
	return existing, nil
}

func (s *cmsService) DeleteReference(refType string, code string) error {
	return s.cmsRepo.DeleteReference(refType, code)
}
