package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type V2Base struct {
	ID        string         `json:"id" gorm:"primaryKey;type:varchar(36);column:id"`
	CreatedAt *time.Time     `json:"created_at" gorm:"column:createdAt;type:timestamp;default:null"`
	UpdatedAt *time.Time     `json:"updated_at" gorm:"column:updatedAt;type:timestamp;default:null"`
	CreatedBy *string        `json:"created_by" gorm:"column:createdBy;type:varchar(255);default:null"`
	UpdatedBy *string        `json:"updated_by" gorm:"column:updatedBy;type:varchar(255);default:null"`
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index;column:deletedAt"`
	IsActive  *bool          `json:"is_active" gorm:"column:isActive;type:boolean;default:true"`
}

func (base *V2Base) BeforeCreate(tx *gorm.DB) (err error) {
	if base.ID == "" {
		base.ID = uuid.New().String()
	}
	if base.IsActive == nil {
		active := true
		base.IsActive = &active
	}
	now := time.Now()
	if base.CreatedAt == nil {
		base.CreatedAt = &now
	}
	if base.UpdatedAt == nil {
		base.UpdatedAt = &now
	}
	return nil
}
