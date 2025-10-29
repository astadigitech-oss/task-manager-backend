package models

import "time"

type TaskImage struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	TaskID     uint      `json:"task_id"`
	URL        string    `json:"url"`
	UploadedBy uint      `json:"uploaded_by"` // Member yg upload
	CreatedAt  time.Time `json:"created_at"`
	Task       Task      `gorm:"foreignKey:TaskID" json:"task"`
}
