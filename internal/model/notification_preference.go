package model

type NotificationPreference struct {
	V2Base
	UserID             uint `json:"-" gorm:"uniqueIndex;not null;column:user_id"`
	DocumentApprovals  bool `json:"document_approvals" gorm:"column:document_approvals;default:true"`
	BrokenLinkReports  bool `json:"broken_link_reports" gorm:"column:broken_link_reports;default:true"`
	SystemUpdates      bool `json:"system_updates" gorm:"column:system_updates;default:false"`
	EmailNotifications bool `json:"email_notifications" gorm:"column:email_notifications;default:true"`
}

func (NotificationPreference) TableName() string {
	return "V2NotificationPreferences"
}

type UpdateNotificationRequest struct {
	DocumentApprovals  *bool `json:"document_approvals" validate:"required"`
	BrokenLinkReports  *bool `json:"broken_link_reports" validate:"required"`
	SystemUpdates      *bool `json:"system_updates" validate:"required"`
	EmailNotifications *bool `json:"email_notifications" validate:"required"`
}
