package repository

import (
	"strings"
	"time"

	"domesv2/config/database"
	"domesv2/internal/model"
	"domesv2/pkg/errors"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type UserRepository interface {
	Create(user *model.User) error
	FindByEmail(email string) (*model.User, error)
	FindByID(id uint) (*model.User, error)
	FindByResetToken(token string) (*model.User, error)
	Update(user *model.User) error
	UpdatePassword(user *model.User, password string) error
	GetNotificationPreferences(userID uint) (*model.NotificationPreference, error)
	UpdateNotificationPreferences(userID uint, prefs *model.NotificationPreference) error
	IsAdminEmail(email string) (bool, error)
	GetAdminEmails() ([]model.AdminEmail, error)
	AddAdminEmail(email string) (*model.AdminEmail, error)
	DeleteAdminEmail(email string) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository() UserRepository {
	return &userRepository{
		db: database.GetDB(),
	}
}

func (r *userRepository) Create(user *model.User) error {
	result := r.db.Create(user)
	if result.Error != nil {
		if strings.Contains(result.Error.Error(), "Duplicate entry") {
			return errors.NewConflictError("User with this email already exists", "USER_ALREADY_EXISTS")
		}
		zap.L().Error("Failed to create user", zap.Error(result.Error))
		return errors.NewInternalServerError("Failed to create user", "DATABASE_ERROR")
	}
	return nil
}

func (r *userRepository) FindByEmail(email string) (*model.User, error) {
	var user model.User
	result := r.db.Where("email = ?", email).First(&user)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, errors.NewNotFoundError("User not found", "USER_NOT_FOUND")
		}
		zap.L().Error("Failed to find user by email", zap.Error(result.Error))
		return nil, errors.NewInternalServerError("Failed to fetch user", "DATABASE_ERROR")
	}
	return &user, nil
}

func (r *userRepository) FindByID(id uint) (*model.User, error) {
	var user model.User
	result := r.db.Preload("NotificationPreferences").First(&user, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, errors.NewNotFoundError("User not found", "USER_NOT_FOUND")
		}
		zap.L().Error("Failed to find user by id", zap.Error(result.Error))
		return nil, errors.NewInternalServerError("Failed to fetch user", "DATABASE_ERROR")
	}
	if user.NotificationPreferences == nil {
		prefs := &model.NotificationPreference{
			UserID:             user.ID,
			DocumentApprovals:  true,
			BrokenLinkReports:  true,
			SystemUpdates:      false,
			EmailNotifications: true,
		}
		if err := r.db.Create(prefs).Error; err != nil {
			zap.L().Error("Failed to create default notification preferences", zap.Error(err))
		} else {
			user.NotificationPreferences = prefs
		}
	}
	return &user, nil
}

func (r *userRepository) FindByResetToken(token string) (*model.User, error) {
	var user model.User
	result := r.db.Where("reset_password_token = ? AND reset_password_expiry > ?", token, time.Now()).First(&user)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, errors.NewBadRequestError("Invalid or expired reset token", "INVALID_RESET_TOKEN")
		}
		zap.L().Error("Failed to find user by reset token", zap.Error(result.Error))
		return nil, errors.NewInternalServerError("Failed to process request", "DATABASE_ERROR")
	}
	return &user, nil
}

func (r *userRepository) Update(user *model.User) error {
	result := r.db.Save(user)
	if result.Error != nil {
		zap.L().Error("Failed to update user", zap.Error(result.Error))
		return errors.NewInternalServerError("Failed to update user", "DATABASE_ERROR")
	}
	return nil
}

func (r *userRepository) UpdatePassword(user *model.User, password string) error {
	result := r.db.Model(user).Updates(map[string]interface{}{
		"password":              password,
		"reset_password_token":  nil,
		"reset_password_expiry": nil,
	})
	if result.Error != nil {
		zap.L().Error("Failed to update password", zap.Error(result.Error))
		return errors.NewInternalServerError("Failed to update password", "DATABASE_ERROR")
	}
	return nil
}

func (r *userRepository) GetNotificationPreferences(userID uint) (*model.NotificationPreference, error) {
	var prefs model.NotificationPreference
	result := r.db.Where("user_id = ?", userID).First(&prefs)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			// Create default
			prefs = model.NotificationPreference{
				UserID:             userID,
				DocumentApprovals:  true,
				BrokenLinkReports:  true,
				SystemUpdates:      false,
				EmailNotifications: true,
			}
			if err := r.db.Create(&prefs).Error; err != nil {
				zap.L().Error("Failed to create default preferences in GetNotificationPreferences", zap.Error(err))
				return nil, errors.NewInternalServerError("Failed to retrieve preferences", "DATABASE_ERROR")
			}
			return &prefs, nil
		}
		zap.L().Error("Failed to fetch preferences", zap.Error(result.Error))
		return nil, errors.NewInternalServerError("Failed to fetch preferences", "DATABASE_ERROR")
	}
	return &prefs, nil
}

func (r *userRepository) UpdateNotificationPreferences(userID uint, prefs *model.NotificationPreference) error {
	result := r.db.Model(&model.NotificationPreference{}).Where("user_id = ?", userID).Updates(map[string]interface{}{
		"document_approvals":  prefs.DocumentApprovals,
		"broken_link_reports": prefs.BrokenLinkReports,
		"system_updates":      prefs.SystemUpdates,
		"email_notifications": prefs.EmailNotifications,
	})
	if result.Error != nil {
		zap.L().Error("Failed to update preferences", zap.Error(result.Error))
		return errors.NewInternalServerError("Failed to update preferences", "DATABASE_ERROR")
	}
	return nil
}

func (r *userRepository) IsAdminEmail(email string) (bool, error) {
	var count int64
	err := r.db.Model(&model.AdminEmail{}).Where("email = ?", email).Count(&count).Error
	if err != nil {
		zap.L().Error("Failed to check admin email whitelist", zap.Error(err))
		return false, errors.NewInternalServerError("Failed to check whitelist status", "DATABASE_ERROR")
	}
	return count > 0, nil
}

func (r *userRepository) GetAdminEmails() ([]model.AdminEmail, error) {
	var emails []model.AdminEmail
	err := r.db.Order("added_at desc").Find(&emails).Error
	if err != nil {
		zap.L().Error("Failed to fetch admin emails", zap.Error(err))
		return nil, errors.NewInternalServerError("Failed to retrieve admin emails", "DATABASE_ERROR")
	}
	return emails, nil
}

func (r *userRepository) AddAdminEmail(email string) (*model.AdminEmail, error) {
	var count int64
	err := r.db.Model(&model.AdminEmail{}).Where("email = ?", email).Count(&count).Error
	if err != nil {
		zap.L().Error("Failed to check existence for new admin email", zap.Error(err))
		return nil, errors.NewInternalServerError("Database query failed", "DATABASE_ERROR")
	}
	if count > 0 {
		return nil, errors.NewConflictError("Admin email already exists", "ADMIN_EMAIL_EXISTS")
	}

	adminEmail := &model.AdminEmail{
		Email:   email,
		AddedAt: time.Now(),
	}
	if err := r.db.Create(adminEmail).Error; err != nil {
		zap.L().Error("Failed to add admin email", zap.Error(err))
		return nil, errors.NewInternalServerError("Failed to save admin email", "DATABASE_ERROR")
	}
	return adminEmail, nil
}

func (r *userRepository) DeleteAdminEmail(email string) error {
	var adminEmail model.AdminEmail
	result := r.db.Where("email = ?", email).First(&adminEmail)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return errors.NewNotFoundError("Admin email not found", "ADMIN_EMAIL_NOT_FOUND")
		}
		zap.L().Error("Failed to find admin email for deletion", zap.Error(result.Error))
		return errors.NewInternalServerError("Database operation failed", "DATABASE_ERROR")
	}

	if err := r.db.Delete(&adminEmail).Error; err != nil {
		zap.L().Error("Failed to delete admin email", zap.Error(err))
		return errors.NewInternalServerError("Failed to delete admin email", "DATABASE_ERROR")
	}
	return nil
}
