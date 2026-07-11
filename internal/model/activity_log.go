package model

type DocumentActivityLog struct {
	V2Base
	DocumentID string    `json:"document_id" gorm:"not null;column:document_id;type:varchar(36)"`
	Document   *Document `json:"document,omitempty" gorm:"foreignKey:DocumentID;constraint:false"`
	Action     string    `json:"action" gorm:"size:50;not null;column:action"` // "view" or "download"
	IPAddress  string    `json:"ip_address" gorm:"size:100;not null;column:ip_address"`
	UserAgent  string    `json:"user_agent" gorm:"size:255;column:user_agent"`
}

func (DocumentActivityLog) TableName() string {
	return "V2DocumentActivityLogs"
}
