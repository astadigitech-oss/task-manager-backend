package models

import (
	"time"

	"gorm.io/gorm"
)

type ProjectImage struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	ProjectID  uint           `json:"project_id"`
	URL        string         `json:"url"`
	UploadedBy uint           `json:"uploaded_by"` // ID Admin uploader
	Project    Project        `gorm:"foreignKey:ProjectID" json:"project"`
	CreatedAt  time.Time      `gorm:"autoCreateTime"`
	UpdatedAt  time.Time      `gorm:"autoUpdateTime"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"` // Soft delete
}
