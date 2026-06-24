package model

import "time"

type AdminEmail struct {
	V2Base
	Email   string    `json:"email" gorm:"uniqueIndex;size:255;column:email"`
	AddedAt time.Time `json:"added_at" gorm:"column:added_at"`
}

func (AdminEmail) TableName() string {
	return "V2AdminEmails"
}

type AddAdminEmailRequest struct {
	Email string `json:"email" validate:"required,email"`
}
