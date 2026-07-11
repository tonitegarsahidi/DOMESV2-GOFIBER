package model

import "time"

type DocumentActivityLog struct {
	ID         uint64     `json:"id" gorm:"primaryKey;autoIncrement;column:id"`
	CreatedAt  *time.Time `json:"created_at" gorm:"column:createdAt;type:timestamp;not null;default:CURRENT_TIMESTAMP"`
	DocumentID string     `json:"document_id" gorm:"not null;column:document_id;type:varchar(36)"`
	Document   *Document  `json:"document,omitempty" gorm:"foreignKey:DocumentID;references:UUID;constraint:false"`
	Action     string     `json:"action" gorm:"size:50;not null;column:action"` // "view" or "download"
	IPAddress  string     `json:"ip_address" gorm:"size:100;not null;column:ip_address"`
	UserAgent  string     `json:"user_agent" gorm:"size:255;column:user_agent"`
}

func (DocumentActivityLog) TableName() string {
	return "V2DocumentActivityLogs"
}
