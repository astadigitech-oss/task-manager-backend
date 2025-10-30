package models

import "time"

type Workspace struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedBy   uint      `json:"created_by"`
	Creator     User      `gorm:"foreignKey:CreatedBy" json:"creator"`
	Projects    []Project `gorm:"foreignKey:WorkspaceID" json:"projects"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`
}
