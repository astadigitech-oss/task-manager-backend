package models

type User struct {
	ID         uint        `gorm:"primaryKey" json:"id"`
	Name       string      `json:"name"`
	Email      string      `json:"email"`
	Password   string      `json:"-"`
	Role       string      `json:"role"`
	Position   *string     `json:"position"` // Jabatan, nullable jika baru ditambah
	Workspaces []Workspace `gorm:"foreignKey:CreatedBy" json:"workspaces"`
	Projects   []Project   `gorm:"many2many:project_users" json:"projects"`
	Tasks      []TaskUser  `gorm:"foreignKey:UserID" json:"tasks"`
}
