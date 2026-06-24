package model

type Report struct {
	V2Base
	DocumentID    string    `json:"document_id" gorm:"not null;column:document_id;type:varchar(36)"`
	Document      *Document `json:"document,omitempty" gorm:"foreignKey:DocumentID;constraint:false"`
	ReporterName  string    `json:"reporter_name" gorm:"size:255;not null;column:reporter_name"`
	ReporterEmail string    `json:"reporter_email" gorm:"size:255;not null;column:reporter_email"`
	Details       string    `json:"details" gorm:"type:text;not null;column:details"`
	Status        string    `json:"status" gorm:"size:50;default:'open';column:status"`
}

func (Report) TableName() string {
	return "V2Reports"
}

type CreateReportRequest struct {
	DocumentID    string `json:"document_id"`
	ReporterName  string `json:"reporter_name"`
	ReporterEmail string `json:"reporter_email"`
	Details       string `json:"details"`
	Captcha       string `json:"captcha,omitempty"`
}

type UpdateReportStatusRequest struct {
	Status string `json:"status"`
}
