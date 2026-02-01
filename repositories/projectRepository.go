// repositories/project_repository.go
package repositories

import (
	"project-management-backend/config"
	"project-management-backend/models"
	"time"

	"gorm.io/gorm"
)

type ProjectRepository interface {
	CreateProject(project *models.Project) error
	GetAllProjects() ([]models.Project, error)
	GetByID(projectID uint) (*models.Project, error)
	GetProjectsByUserID(userID uint) ([]models.Project, error)
	GetProjectsByWorkspace(workspaceID uint) ([]models.Project, error)
	UpdateProject(project *models.Project) error
	SoftDeleteProject(projectID uint) error
	DeleteProject(projectID uint) error
	AddMember(pu *models.ProjectUser) error
	GetMembers(projectID uint) ([]models.ProjectUser, error)
	GetUserByID(userID uint) (*models.User, error)
	GetProjectMember(projectID uint, userID uint) (*models.ProjectUser, error)
	IsUserMember(projectID uint, userID uint) (bool, error)
	IsUserInWorkspace(workspaceID uint, userID uint) (bool, error)
	RemoveMember(projectID uint, userID uint) error
	RemoveMembers(projectID uint, userIDs []uint) error
	GetActivityLogsSince(projectID uint, since time.Time) ([]models.ActivityLog, error)
}

type projectRepository struct{}

func NewProjectRepository() ProjectRepository {
	return &projectRepository{}
}

func (r *projectRepository) GetActivityLogsSince(projectID uint, since time.Time) ([]models.ActivityLog, error) {
	var activities []models.ActivityLog
	db := config.DB

	var taskIDs []uint
	if err := db.Model(&models.Task{}).Where("project_id = ?", projectID).Pluck("id", &taskIDs).Error; err != nil {
		return nil, err
	}

	if len(taskIDs) == 0 {
		return []models.ActivityLog{}, nil
	}

	err := db.Where("table_name = 'tasks' AND item_id IN (?) AND created_at >= ?", taskIDs, since).Find(&activities).Error
	return activities, err
}

func (r *projectRepository) CreateProject(project *models.Project) error {
	return config.DB.Create(project).Error
}

func (r *projectRepository) GetAllProjects() ([]models.Project, error) {
	var projects []models.Project
	err := config.DB.
		Joins("JOIN workspaces ON workspaces.id = projects.workspace_id").
		Where("projects.deleted_at IS NULL AND workspaces.deleted_at IS NULL").
		Preload("Members.User").
		Preload("Tasks", func(db *gorm.DB) *gorm.DB {
			return db.Where("deleted_at IS NULL")
		}).
		Preload("Images", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, url, project_id")
		}).
		Find(&projects).Error
	return projects, err
}

func (r *projectRepository) GetProjectsByUserID(userID uint) ([]models.Project, error) {
	var projects []models.Project
	err := config.DB.
		Select("DISTINCT projects.*").
		Joins("JOIN project_users ON project_users.project_id = projects.id").
		Joins("JOIN workspaces ON workspaces.id = projects.workspace_id").
		Where("project_users.user_id = ? AND projects.deleted_at IS NULL AND workspaces.deleted_at IS NULL", userID).
		Preload("Members.User").
		Preload("Tasks", func(db *gorm.DB) *gorm.DB {
			return db.Where("deleted_at IS NULL")
		}).
		Preload("Images", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, url, project_id")
		}).
		Find(&projects).Error
	return projects, err
}

func (r *projectRepository) GetProjectsByWorkspace(workspaceID uint) ([]models.Project, error) {
	var projects []models.Project
	err := config.DB.
		Joins("JOIN workspaces ON workspaces.id = projects.workspace_id").
		Where("projects.workspace_id = ? AND projects.deleted_at IS NULL AND workspaces.deleted_at IS NULL", workspaceID).
		Preload("Members.User").
		Preload("Tasks", func(db *gorm.DB) *gorm.DB {
			return db.Where("deleted_at IS NULL")
		}).
		Preload("Images", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, url, project_id")
		}).
		Find(&projects).Error
	return projects, err
}

func (r *projectRepository) GetByID(projectID uint) (*models.Project, error) {
	var project models.Project
	err := config.DB.
		Joins("JOIN workspaces ON workspaces.id = projects.workspace_id").
		Where("projects.id = ? AND projects.deleted_at IS NULL AND workspaces.deleted_at IS NULL", projectID).
		Preload("Members.User").
		Preload("Tasks", "deleted_at IS NULL").
		Preload("Images").
		Preload("Workspace").
		First(&project).Error
	return &project, err
}

func (r *projectRepository) UpdateProject(project *models.Project) error {
	return config.DB.Model(&models.Project{}).
		Where("id = ? AND deleted_at IS NULL", project.ID).
		Updates(map[string]interface{}{
			"name":        project.Name,
			"description": project.Description,
			"updated_at":  time.Now(),
		}).Error
}

func (r *projectRepository) SoftDeleteProject(projectID uint) error {
	return config.DB.Model(&models.Project{}).
		Where("id = ?", projectID).
		Update("deleted_at", time.Now()).Error
}

func (r *projectRepository) DeleteProject(projectID uint) error {
	return config.DB.Transaction(func(tx *gorm.DB) error {
		// 1. Delete task_users yang terkait
		if err := tx.Exec(`
			DELETE tu FROM task_users tu
			INNER JOIN tasks t ON t.id = tu.task_id
			WHERE t.project_id = ?
		`, projectID).Error; err != nil {
			return err
		}

		// 2. Delete tasks yang terkait
		if err := tx.Where("project_id = ?", projectID).Delete(&models.Task{}).Error; err != nil {
			return err
		}

		// 3. Delete project_images yang terkait
		if err := tx.Where("project_id = ?", projectID).Delete(&models.ProjectImage{}).Error; err != nil {
			return err
		}

		// 4. Delete project_users (members)
		if err := tx.Where("project_id = ?", projectID).Delete(&models.ProjectUser{}).Error; err != nil {
			return err
		}

		// 5. HARD DELETE project
		if err := tx.Unscoped().Where("id = ?", projectID).Delete(&models.Project{}).Error; err != nil {
			return err
		}

		return nil
	})
}

func (r *projectRepository) AddMember(pu *models.ProjectUser) error {
	return config.DB.Create(pu).Error
}

func (r *projectRepository) GetMembers(projectID uint) ([]models.ProjectUser, error) {
	var members []models.ProjectUser
	err := config.DB.
		Preload("User").
		Where("project_id = ?", projectID).
		Find(&members).Error
	return members, err
}

func (r *projectRepository) GetProjectMember(projectID uint, userID uint) (*models.ProjectUser, error) {
	var projectUser models.ProjectUser
	err := config.DB.
		Where("project_id = ? AND user_id = ?", projectID, userID).
		First(&projectUser).Error
	if err != nil {
		return nil, err
	}
	return &projectUser, nil
}

func (r *projectRepository) GetUserByID(userID uint) (*models.User, error) {
	var user models.User
	err := config.DB.First(&user, userID).Error
	return &user, err
}

func (r *projectRepository) IsUserMember(projectID uint, userID uint) (bool, error) {
	var count int64
	err := config.DB.Model(&models.ProjectUser{}).
		Where("project_id = ? AND user_id = ?", projectID, userID).
		Count(&count).Error
	return count > 0, err
}

func (r *projectRepository) IsUserInWorkspace(workspaceID uint, userID uint) (bool, error) {
	var count int64
	err := config.DB.Model(&models.WorkspaceUser{}).
		Where("workspace_id = ? AND user_id = ?", workspaceID, userID).
		Count(&count).Error
	return count > 0, err
}

func (r *projectRepository) RemoveMember(projectID uint, userID uint) error {
	return config.DB.
		Where("project_id = ? AND user_id = ?", projectID, userID).
		Delete(&models.ProjectUser{}).Error
}

func (r *projectRepository) RemoveMembers(projectID uint, userIDs []uint) error {
	return config.DB.
		Where("project_id = ? AND user_id IN ?", projectID, userIDs).
		Delete(&models.ProjectUser{}).Error
}
