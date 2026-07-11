package model

import "time"

type DocumentStats struct {
	DocumentID         uint      `json:"document_id" gorm:"primaryKey;column:document_id"`
	Document           *Document `json:"document,omitempty" gorm:"foreignKey:DocumentID;constraint:false"`
	TotalViews         int       `json:"total_views" gorm:"default:0;column:total_views"`
	TotalDownloads     int       `json:"total_downloads" gorm:"default:0;column:total_downloads"`
	LastProcessedLogID uint64    `json:"last_processed_log_id" gorm:"default:0;column:last_processed_log_id"`
	LastUpdate         time.Time `json:"last_update" gorm:"column:last_update;type:timestamp;not null;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"`
}

func (DocumentStats) TableName() string {
	return "V2DocumentStats"
}
