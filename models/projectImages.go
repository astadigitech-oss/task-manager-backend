package models

import "time"

type ProjectImage struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	ProjectID  uint      `json:"project_id"`
	URL        string    `json:"url"`
	UploadedBy uint      `json:"uploaded_by"` // ID Admin uploader
	CreatedAt  time.Time `json:"created_at"`
	Project    Project   `gorm:"foreignKey:ProjectID" json:"project"`
}
