// controllers/project_controller.go
package controllers

import (
	"project-management-backend/models"
	"project-management-backend/services"
	"project-management-backend/utils"

	"github.com/gin-gonic/gin"
)

type ProjectController struct {
	Service services.ProjectService
}

func NewProjectController(service services.ProjectService) *ProjectController {
	return &ProjectController{Service: service}
}

func (pc *ProjectController) ListProjects(c *gin.Context) {
	currentUser := GetCurrentUser(c)

	projects, err := pc.Service.GetAllProjects(currentUser)
	if err != nil {
		c.JSON(403, gin.H{"error": err.Error()})
		return
	}

	var projectList []gin.H
	for _, project := range projects {
		projectList = append(projectList, gin.H{
			"id":           project.ID,
			"name":         project.Name,
			"description":  project.Description,
			"workspace_id": project.WorkspaceID,
			"member_count": len(project.Members),
			"task_count":   len(project.Tasks),
			"image_count":  len(project.Images),
		})
	}

	c.JSON(200, APIResponse{
		Success: true,
		Code:    200,
		Message: "List project berhasil di ambil",
		Data:    projectList,
	})
}

func (pc *ProjectController) CreateProject(c *gin.Context) {
	var input struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		WorkspaceID uint   `json:"workspace_id"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	currentUser := GetCurrentUser(c)

	project := models.Project{
		Name:        input.Name,
		Description: input.Description,
		WorkspaceID: input.WorkspaceID,
	}

	if err := pc.Service.CreateProject(&project, currentUser); err != nil {
		c.JSON(403, gin.H{"error": err.Error()})
		return
	}

	utils.ActivityLog(currentUser.ID, "CREATE_PROJECT", "project", project.ID, nil, project)

	c.JSON(201, APIResponse{
		Success: true,
		Code:    201,
		Message: "Project berhasil di buat",
		Data: gin.H{
			"id":           project.ID,
			"name":         project.Name,
			"description":  project.Description,
			"workspace_id": project.WorkspaceID,
		},
	})
}

func (pc *ProjectController) DetailProject(c *gin.Context) {
	projectID, err := ParseUintParam(c, "project_id")
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	currentUser := GetCurrentUser(c)

	project, err := pc.Service.GetByID(projectID, currentUser)
	if err != nil {
		c.JSON(403, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, APIResponse{
		Success: true,
		Code:    200,
		Message: "Detail project berhasil diambil",
		Data: gin.H{
			"id":           project.ID,
			"name":         project.Name,
			"description":  project.Description,
			"workspace_id": project.WorkspaceID,
			"created_by":   project.CreatedBy,
		},
	})
}

func (pc *ProjectController) UpdateProject(c *gin.Context) {
	projectID, err := ParseUintParam(c, "project_id")
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	var input struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	currentUser := GetCurrentUser(c)

	// Get data lama sebelum update
	oldProject, err := pc.Service.GetByID(projectID, currentUser)
	if err != nil {
		c.JSON(403, gin.H{"error": "Project tidak ditemukan"})
		return
	}

	project := models.Project{
		ID:          projectID,
		Name:        input.Name,
		Description: input.Description,
	}

	if err := pc.Service.UpdateProject(&project, currentUser); err != nil {
		c.JSON(403, gin.H{"error": err.Error()})
		return
	}

	// Get data baru setelah update
	updatedProject, err := pc.Service.GetByID(projectID, currentUser)
	if err != nil {
		c.JSON(403, gin.H{"error": "Gagal mengambil data project setelah update"})
		return
	}

	utils.ActivityLog(currentUser.ID, "UPDATE_PROJECT", "project", projectID, oldProject, updatedProject)

	c.JSON(200, APIResponse{
		Success: true,
		Code:    200,
		Message: "Project berhasil di update",
		Data: gin.H{
			"id":           updatedProject.ID,
			"name":         updatedProject.Name,
			"description":  updatedProject.Description,
			"workspace_id": updatedProject.WorkspaceID,
		},
	})
}

func (pc *ProjectController) SoftDeleteProject(c *gin.Context) {
	projectID, err := ParseUintParam(c, "project_id")
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	currentUser := GetCurrentUser(c)

	oldProject, err := pc.Service.GetByID(projectID, currentUser)
	if err != nil {
		c.JSON(403, gin.H{"error": "Project tidak ditemukan"})
		return
	}

	if err := pc.Service.SoftDeleteProject(projectID, currentUser); err != nil {
		c.JSON(403, gin.H{"error": err.Error()})
		return
	}

	utils.ActivityLog(currentUser.ID, "SOFT_DELETE_PROJECT", "project", projectID, oldProject, "deleted")
	c.JSON(200, APIResponse{
		Success: true,
		Code:    200,
		Message: "Project berhasil di soft delete",
		Data:    gin.H{"project_id": projectID},
	})
}

func (pc *ProjectController) DeleteProject(c *gin.Context) {
	projectID, err := ParseUintParam(c, "project_id")
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	currentUser := GetCurrentUser(c)

	project, err := pc.Service.GetByID(projectID, currentUser)
	if err != nil {
		c.JSON(404, gin.H{"error": "Project tidak ditemukan"})
		return
	}

	var input struct {
		Confirm bool `json:"confirm"`
	}
	if err := c.ShouldBindJSON(&input); err != nil || !input.Confirm {
		c.JSON(400, gin.H{
			"error":   "Konfirmasi diperlukan untuk hard delete",
			"warning": "Tindakan ini akan menghapus PERMANEN:",
			"data_akan_dihapus": gin.H{
				"project":       project.Name,
				"tasks_count":   len(project.Tasks),
				"members_count": len(project.Members),
				"images_count":  len(project.Images),
			},
			"confirmation_required": true,
		})
		return
	}

	if err := pc.Service.DeleteProject(projectID, currentUser); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	utils.ActivityLog(currentUser.ID, "DELETE_PROJECT", "project", projectID, project, "hard deleted")

	c.JSON(200, APIResponse{
		Success: true,
		Code:    200,
		Message: "Project dan semua data terkait berhasil dihapus permanen",
		Data: gin.H{
			"project_id":      projectID,
			"project_name":    project.Name,
			"deleted_tasks":   len(project.Tasks),
			"deleted_members": len(project.Members),
			"deleted_images":  len(project.Images),
		},
	})
}

func (pc *ProjectController) AddMember(c *gin.Context) {
	projectID, err := ParseUintParam(c, "project_id")
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	var input struct {
		UserID uint   `json:"user_id"`
		Role   string `json:"role_in_project"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	currentUser := GetCurrentUser(c)

	if err := pc.Service.AddMember(projectID, input.UserID, input.Role, currentUser); err != nil {
		c.JSON(403, gin.H{"error": err.Error()})
		return
	}

	utils.ActivityLog(currentUser.ID, "ADD_MEMBER_PROJECT", "project", projectID, nil, input)

	c.JSON(200, APIResponse{
		Success: true,
		Code:    200,
		Message: "Member berhasil di buat",
		Data: gin.H{
			"project_id": projectID,
			"user_id":    input.UserID,
			"role":       input.Role,
		},
	})
}

func (pc *ProjectController) GetMembers(c *gin.Context) {
	projectID, err := ParseUintParam(c, "project_id")
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	currentUser := GetCurrentUser(c)

	members, err := pc.Service.GetMembers(projectID, currentUser)
	if err != nil {
		c.JSON(403, gin.H{"error": err.Error()})
		return
	}

	memberProject := make([]gin.H, 0)
	for _, member := range members {
		memberProject = append(memberProject, gin.H{
			"id":         member.User.ID,
			"name":       member.User.Name,
			"user_email": member.User.Email,
			"role":       member.User.Role,
		})
	}

	c.JSON(200, APIResponse{
		Success: true,
		Code:    200,
		Message: "Members project berhasil diambil",
		Data:    memberProject,
	})
}
