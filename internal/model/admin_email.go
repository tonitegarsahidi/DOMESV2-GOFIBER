package model

import "time"

type AdminEmail struct {
	Email   string    `json:"email" gorm:"primaryKey;size:255;column:email"`
	AddedAt time.Time `json:"added_at" gorm:"column:added_at"`
}

func (AdminEmail) TableName() string {
	return "AdminEmails"
}

type AddAdminEmailRequest struct {
	Email string `json:"email" validate:"required,email"`
}
