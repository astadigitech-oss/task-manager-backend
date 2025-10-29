package models

import (
	"time"
)

type Task struct {
	ID          uint        `gorm:"primaryKey" json:"id"`
	ProjectID   uint        `json:"project_id"`
	Title       string      `json:"title"`
	Description string      `json:"description"`
	Status      string      `json:"status"`
	Priority    string      `json:"priority"`
	StartDate   time.Time   `json:"start_date"`
	DueDate     time.Time   `json:"due_date"`
	Notes       *string     `json:"notes"` // Notes yang diisi member, nullable
	Project     Project     `gorm:"foreignKey:ProjectID" json:"project"`
	Members     []TaskUser  `gorm:"foreignKey:TaskID" json:"members"`
	Images      []TaskImage `gorm:"foreignKey:TaskID" json:"images"` // Tambah relasi image di task
}
