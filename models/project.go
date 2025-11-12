package models

import (
	"time"

	"gorm.io/gorm"
)

type Project struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	CreatedBy   uint           `json:"created_by"`
	WorkspaceID uint           `json:"workspace_id"`
	Workspace   Workspace      `gorm:"foreignKey:WorkspaceID;constraint:OnDelete:CASCADE" json:"workspace"`
	Members     []ProjectUser  `gorm:"foreignKey:ProjectID;constraint:OnDelete:CASCADE" json:"members"`
	Tasks       []Task         `gorm:"foreignKey:ProjectID;constraint:OnDelete:CASCADE" json:"tasks"`
	Images      []ProjectImage `gorm:"foreignKey:ProjectID;constraint:OnDelete:CASCADE" json:"images"`
	CreatedAt   time.Time      `gorm:"autoCreateTime"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"` // Soft delete
}
