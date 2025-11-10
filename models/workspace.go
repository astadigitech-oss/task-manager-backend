package models

import (
	"time"

	"gorm.io/gorm"
)

type Workspace struct {
	ID          uint            `gorm:"primaryKey" json:"id"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	CreatedBy   uint            `json:"created_by"`
	CreatedAt   time.Time       `gorm:"autoCreateTime"`
	UpdatedAt   time.Time       `gorm:"autoUpdateTime"`
	DeletedAt   gorm.DeletedAt  `gorm:"index" json:"deleted_at,omitempty"`
	Members     []WorkspaceUser `gorm:"foreignKey:WorkspaceID" json:"members"` // keaggotaan lewat pivot
	Projects    []Project       `gorm:"foreignKey:WorkspaceID" json:"projects"`
}
