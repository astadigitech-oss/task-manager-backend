package models

type Workspace struct {
	ID          uint            `gorm:"primaryKey" json:"id"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	CreatedBy   uint            `json:"created_by"`                            // owner``
	Members     []WorkspaceUser `gorm:"foreignKey:WorkspaceID" json:"members"` // keaggotaan lewat pivot
	Projects    []Project       `gorm:"foreignKey:WorkspaceID" json:"projects"`
}
