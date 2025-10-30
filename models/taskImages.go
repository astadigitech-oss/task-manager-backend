package models

import "time"

type TaskImage struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	TaskID     uint      `json:"task_id"`
	URL        string    `json:"url"`
	UploadedBy uint      `json:"uploaded_by"` // Member yg upload
	Task       Task      `gorm:"foreignKey:TaskID" json:"task"`
	CreatedAt  time.Time `gorm:"autoCreateTime"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime"`
}
