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
	result := r.db.First(&user, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, errors.NewNotFoundError("User not found", "USER_NOT_FOUND")
		}
		zap.L().Error("Failed to find user by id", zap.Error(result.Error))
		return nil, errors.NewInternalServerError("Failed to fetch user", "DATABASE_ERROR")
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
