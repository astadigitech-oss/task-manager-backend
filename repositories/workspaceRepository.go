package repositories

import (
	"project-management-backend/config"
	"project-management-backend/models"
	"time"

	"gorm.io/gorm"
)

type WorkspaceRepository interface {
	CreateWorkspace(workspace *models.Workspace) error
	GetAllWorkspaces() ([]models.Workspace, error)
	UpdateWorkspace(workspace *models.Workspace) error
	SoftDeleteWorkspace(workspaceID uint) error
	DeleteWorkspace(workspaceID uint) error // Hard delete
	AddMember(wu *models.WorkspaceUser) error
	GetMembers(workspaceID uint) ([]models.WorkspaceUser, error)
	GetByID(workspaceID uint) (*models.Workspace, error)
	GetUserByID(userID uint) (*models.User, error)
	IsUserMember(workspaceID uint, userID uint) (bool, error)
	GetWorkspaceMember(workspaceID uint, userID uint) (*models.WorkspaceUser, error)
}

type workspaceRepository struct{}

func NewWorkspaceRepository() WorkspaceRepository {
	return &workspaceRepository{}
}

func (r *workspaceRepository) CreateWorkspace(workspace *models.Workspace) error {
	return config.DB.Create(workspace).Error
}

func (r *workspaceRepository) GetAllWorkspaces() ([]models.Workspace, error) {
	var workspaces []models.Workspace
	err := config.DB.
		Where("deleted_at IS NULL").
		Preload("Projects").
		Preload("Members").
		Find(&workspaces).Error
	return workspaces, err
}

func (r *workspaceRepository) GetByID(workspaceID uint) (*models.Workspace, error) {
	var workspace models.Workspace
	err := config.DB.
		Where("id = ? AND deleted_at IS NULL", workspaceID).
		Preload("Members.User").
		Preload("Projects").
		First(&workspace, workspaceID).Error
	return &workspace, err
}

func (r *workspaceRepository) UpdateWorkspace(workspace *models.Workspace) error {
	return config.DB.Model(&models.Workspace{}).
		Where("id = ? AND deleted_at IS NULL", workspace.ID).
		Updates(map[string]interface{}{
			"name":        workspace.Name,
			"description": workspace.Description,
			"updated_at":  time.Now(),
		}).Error
}

func (r *workspaceRepository) SoftDeleteWorkspace(workspaceID uint) error {
	return config.DB.Model(&models.Workspace{}).
		Where("id = ?", workspaceID).
		Update("deleted_at", time.Now()).Error
}

func (r *workspaceRepository) DeleteWorkspace(workspaceID uint) error {
	return config.DB.Transaction(func(tx *gorm.DB) error {
		// 1. Delete task_users yang terkait
		if err := tx.Exec(`
            DELETE tu FROM task_users tu
            INNER JOIN tasks t ON t.id = tu.task_id
            INNER JOIN projects p ON p.id = t.project_id
            WHERE p.workspace_id = ?
        `, workspaceID).Error; err != nil {
			return err
		}

		// 2. Delete tasks yang terkait
		if err := tx.Exec(`
            DELETE t FROM tasks t
            INNER JOIN projects p ON p.id = t.project_id
            WHERE p.workspace_id = ?
        `, workspaceID).Error; err != nil {
			return err
		}

		// 3. Delete project_users yang terkait
		if err := tx.Exec(`
            DELETE pu FROM project_users pu
            INNER JOIN projects p ON p.id = pu.project_id
            WHERE p.workspace_id = ?
        `, workspaceID).Error; err != nil {
			return err
		}

		// 4. Delete projects yang terkait
		if err := tx.Where("workspace_id = ?", workspaceID).Delete(&models.Project{}).Error; err != nil {
			return err
		}

		// 5. Delete workspace_users (members)
		if err := tx.Where("workspace_id = ?", workspaceID).Delete(&models.WorkspaceUser{}).Error; err != nil {
			return err
		}

		// 6. Sekarang baru delete workspace
		if err := tx.Unscoped().Where("id = ?", workspaceID).Delete(&models.Workspace{}).Error; err != nil {
			return err
		}

		return nil
	})
}

func (r *workspaceRepository) AddMember(wu *models.WorkspaceUser) error {
	return config.DB.Create(wu).Error
}

func (r *workspaceRepository) GetMembers(workspaceID uint) ([]models.WorkspaceUser, error) {
	var members []models.WorkspaceUser
	err := config.DB.Preload("User").Where("workspace_id = ?", workspaceID).Find(&members).Error
	return members, err
}

func (r *workspaceRepository) GetUserByID(userID uint) (*models.User, error) {
	var user models.User
	err := config.DB.First(&user, userID).Error
	return &user, err
}

func (r *workspaceRepository) IsUserMember(workspaceID uint, userID uint) (bool, error) {
	var count int64
	err := config.DB.Model(&models.WorkspaceUser{}).
		Where("workspace_id = ? AND user_id = ?", workspaceID, userID).
		Count(&count).Error
	return count > 0, err
}

func (r *workspaceRepository) GetWorkspaceMember(workspaceID uint, userID uint) (*models.WorkspaceUser, error) {
	var workspaceUser models.WorkspaceUser
	err := config.DB.Where("workspace_id = ? AND user_id = ?", workspaceID, userID).First(&workspaceUser).Error
	if err != nil {
		return nil, err
	}
	return &workspaceUser, nil
}
