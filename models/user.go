package models

import "time"

type User struct {
	ID         uint        `gorm:"primaryKey" json:"id"`
	Name       string      `json:"name"`
	Email      string      `gorm:"unique;not null" json:"email"`
	Password   string      `json:"-"`
	Role       string      `json:"role"`
	Position   *string     `json:"position"` // Jabatan, nullable
	Workspaces []Workspace `gorm:"many2many:workspace_users" json:"workspaces"`
	Projects   []Project   `gorm:"many2many:project_users" json:"projects"`
	Tasks      []TaskUser  `gorm:"foreignKey:UserID" json:"tasks"`
	CreatedAt  time.Time   `gorm:"autoCreateTime"`
	UpdatedAt  time.Time   `gorm:"autoUpdateTime"`
}
