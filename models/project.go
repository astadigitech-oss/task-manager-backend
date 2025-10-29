package models

type Project struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	CreatedBy   uint           `json:"created_by"`
	WorkspaceID uint           `json:"workspace_id"`
	Workspace   Workspace      `gorm:"foreignKey:WorkspaceID" json:"workspace"`
	Members     []ProjectUser  `gorm:"foreignKey:ProjectID" json:"members"`
	Tasks       []Task         `gorm:"foreignKey:ProjectID" json:"tasks"`
	Images      []ProjectImage `gorm:"foreignKey:ProjectID" json:"images"` // Tambah relasi image di project
}
