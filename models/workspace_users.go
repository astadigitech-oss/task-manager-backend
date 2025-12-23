package models

import "time"

type WorkspaceUser struct {
	ID              uint      `gorm:"primaryKey" json:"id"`
	WorkspaceID     uint      `json:"workspace_id"`
	UserID          uint      `json:"user_id"`
	RoleInWorkspace *string   `json:"role_in_workspace"` // granular
	Workspace       Workspace `gorm:"foreignKey:WorkspaceID;constraint:OnDelete:CASCADE" json:"workspaces"`
	User            User      `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"user"`
	CreatedAt       time.Time `gorm:"autoCreateTime"`
	UpdatedAt       time.Time `gorm:"autoUpdateTime"`
}

func (WorkspaceUser) TableName() string { return "workspace_users" }
