package models

type WorkspaceUser struct {
	ID              uint      `gorm:"primaryKey" json:"id"`
	WorkspaceID     uint      `json:"workspace_id"`
	UserID          uint      `json:"user_id"`
	RoleInWorkspace *string   `json:"role_in_workspace"` // granular
	Workspace       Workspace `gorm:"foreignKey:WorkspaceID" json:"workspaces"`
	User            User      `gorm:"foreignKey:UserID" json:"user"`
}

func (WorkspaceUser) TableName() string { return "workspace_users" }
