package model

import (
	"time"
)

type User struct {
	ID                      uint                    `json:"id" gorm:"primaryKey;column:id"`
	Username                *string                 `json:"username" gorm:"uniqueIndex:username;size:255"`
	Name                    *string                 `json:"name" gorm:"size:255"`
	FirstName               *string                 `json:"first_name" gorm:"column:first_name;size:255"`
	LastName                *string                 `json:"last_name" gorm:"column:last_name;size:255"`
	Password                string                  `json:"-" gorm:"size:255"`
	Type                    *string                 `json:"type" gorm:"size:255"`
	Role                    *string                 `json:"role" gorm:"size:50"`
	Status                  *string                 `json:"status" gorm:"size:20;default:'active'"`
	Position                *string                 `json:"position" gorm:"size:255"`
	Organization            *string                 `json:"organization" gorm:"size:255"`
	PhoneNumber             *string                 `json:"phone_number" gorm:"column:phone_number;size:255"`
	Email                   string                  `json:"email" gorm:"unique;size:255;not null"`
	AvatarURL               *string                 `json:"avatar_url" gorm:"column:avatar_url;size:255"`
	RegistrationID          *string                 `json:"registration_id" gorm:"column:registration_id;type:char(36)"`
	Metadata                *string                 `json:"metadata" gorm:"type:json"`
	ResetPasswordToken      *string                 `json:"-" gorm:"column:reset_password_token;size:255"`
	ResetPasswordExpiry     *time.Time              `json:"-" gorm:"column:reset_password_expiry"`
	NotificationPreferences *NotificationPreference `json:"notification_preferences,omitempty" gorm:"foreignKey:UserID;constraint:false"`
	CreatedAt               time.Time               `json:"created_at" gorm:"column:createdAt"`
	UpdatedAt               time.Time               `json:"updated_at" gorm:"column:updatedAt"`
}

func (User) TableName() string {
	return "Users"
}

type RegisterRequest struct {
	FirstName       string `json:"first_name"`
	LastName        string `json:"last_name"`
	Position        string `json:"position"`
	Organization    string `json:"organization"`
	PhoneNumber     string `json:"phone_number"`
	Email           string `json:"email"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
	Captcha         string `json:"captcha,omitempty"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Captcha  string `json:"captcha,omitempty"`
}

type ForgotPasswordRequest struct {
	Email   string `json:"email"`
	Captcha string `json:"captcha,omitempty"`
}

type ResetPasswordRequest struct {
	Token           string `json:"token"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
}

type AuthResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

type UserProfileResponse struct {
	ID                      uint                    `json:"id"`
	Username                *string                 `json:"username"`
	Name                    *string                 `json:"name"`
	FirstName               *string                 `json:"first_name"`
	LastName                *string                 `json:"last_name"`
	Email                   string                  `json:"email"`
	Type                    *string                 `json:"type"`
	Role                    *string                 `json:"role"`
	Status                  *string                 `json:"status"`
	Position                *string                 `json:"position"`
	Organization            *string                 `json:"organization"`
	PhoneNumber             *string                 `json:"phone_number"`
	AvatarURL               *string                 `json:"avatar_url"`
	NotificationPreferences *NotificationPreference `json:"notification_preferences,omitempty"`
	CreatedAt               time.Time               `json:"created_at"`
	UpdatedAt               time.Time               `json:"updated_at"`
}

type UpdateProfileRequest struct {
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	Position     string `json:"position"`
	Organization string `json:"organization"`
	PhoneNumber  string `json:"phone_number"`
}

type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password"`
	NewPassword     string `json:"new_password"`
	ConfirmPassword string `json:"confirm_password"`
}

type CreateUserRequest struct {
	FirstName       string `json:"first_name"`
	LastName        string `json:"last_name"`
	Email           string `json:"email"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
	Organization    string `json:"organization"`
	Position        string `json:"position"`
	PhoneNumber     string `json:"phone_number"`
	Role            string `json:"role"`
	Status          string `json:"status"`
}

type UpdateUserRequest struct {
	FirstName    *string `json:"first_name"`
	LastName     *string `json:"last_name"`
	Organization *string `json:"organization"`
	Position     *string `json:"position"`
	PhoneNumber  *string `json:"phone_number"`
	Role         *string `json:"role"`
	Status       *string `json:"status"`
}
