package models

import (
	"time"

	"gorm.io/gorm"
)

type Task struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	ProjectID   uint      `json:"project_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	Priority    string    `json:"priority"`
	StartDate   time.Time `json:"start_date"`
	DueDate     time.Time `json:"due_date"`
	Notes       *string   `json:"notes"` // Notes yang diisi member, nullable

	Project   Project        `gorm:"foreignKey:ProjectID;constraint:OnDelete:CASCADE" json:"project"`
	Members   []TaskUser     `gorm:"foreignKey:TaskID;constraint:OnDelete:CASCADE" json:"members"`
	Images    []TaskImage    `gorm:"foreignKey:TaskID;constraint:OnDelete:CASCADE" json:"images"`
	CreatedAt time.Time      `gorm:"autoCreateTime"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"` // Soft delete
}
